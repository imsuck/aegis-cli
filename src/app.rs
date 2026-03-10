use crate::vault::Entry;
use crossterm::event::{self, Event, KeyCode, KeyEventKind};
use nucleo::pattern::{CaseMatching, Normalization, Pattern};
use nucleo::{Config, Matcher, Utf32Str};
use ratatui::Terminal;
use ratatui::backend::CrosstermBackend;
use std::io::Stdout;
use zeroize::Zeroizing;

#[derive(Debug)]
pub struct App {
    pub entries: Vec<Entry>,
    pub selected_index: usize,
    pub password: Zeroizing<String>,
    pub search_query: String,
    pub search_mode: bool,
    pub show_code: bool,
    exit: bool,
}

impl App {
    pub fn new() -> Self {
        Self {
            entries: Vec::new(),
            selected_index: 0,
            password: Zeroizing::new(String::new()),
            search_query: String::new(),
            search_mode: false,
            show_code: false,
            exit: false,
        }
    }

    pub fn set_entries(&mut self, entries: Vec<Entry>) {
        self.entries = entries;
    }

    pub fn set_password(&mut self, password: String) {
        self.password = Zeroizing::new(password);
    }

    /// Run the main application loop
    pub fn run(&mut self, terminal: &mut Terminal<CrosstermBackend<Stdout>>) -> io::Result<()> {
        while !self.exit {
            terminal.draw(|frame| self.draw(frame))?;
            self.handle_events()?;
        }
        Ok(())
    }

    fn handle_events(&mut self) -> io::Result<()> {
        if event::poll(std::time::Duration::from_millis(100))? {
            if let Event::Key(key) = event::read()? {
                if key.kind == KeyEventKind::Press {
                    self.handle_key_event(key.code);
                }
            }
        }
        Ok(())
    }

    fn handle_key_event(&mut self, key: KeyCode) {
        match key {
            KeyCode::Char('q') => self.exit(),
            KeyCode::Char('j') | KeyCode::Down => self.move_selection(1),
            KeyCode::Char('k') | KeyCode::Up => self.move_selection(-1),
            KeyCode::Char('/') => self.enter_search_mode(),
            KeyCode::Esc if self.search_mode => self.exit_search_mode(),
            KeyCode::Enter if self.search_mode => self.exit_search_mode(),
            KeyCode::Char(c) if self.search_mode => self.search_query.push(c),
            KeyCode::Char('c') => self.toggle_show_code(),
            KeyCode::Char('y') if self.show_code => self.yank_code(),
            _ => {}
        }
    }

    fn exit(&mut self) {
        self.exit = true;
    }

    fn move_selection(&mut self, delta: isize) {
        let filtered = self.filtered_entries();
        if filtered.is_empty() {
            return;
        }
        let new_index = (self.selected_index as isize + delta)
            .clamp(0, filtered.len() as isize - 1) as usize;
        self.selected_index = new_index;
    }

    fn enter_search_mode(&mut self) {
        self.search_mode = true;
        self.search_query.clear();
    }

    fn exit_search_mode(&mut self) {
        self.search_mode = false;
    }

    fn toggle_show_code(&mut self) {
        self.show_code = !self.show_code;
    }

    fn yank_code(&mut self) {
        if let Some(code) = self.yank_current_code() {
            if let Ok(mut clipboard) = arboard::Clipboard::new() {
                let _ = clipboard.set_text(&code);
            }
        }
    }

    pub fn filtered_entries(&self) -> Vec<&Entry> {
        if self.search_query.is_empty() {
            return self.entries.iter().collect();
        }

        // Check for property prefix: %<property> <query>
        if self.search_query.starts_with('%') {
            let parts: Vec<&str> = self.search_query[1..].splitn(2, ' ').collect();
            if parts.len() == 2 {
                let prop = parts[0].to_lowercase();
                let query = parts[1];
                return self.filter_by_property(&prop, query);
            }
        }

        // Default: fuzzy search in issuer field using nucleo
        self.fuzzy_filter_entries(&self.search_query, |e| e.issuer.as_str())
    }

    pub fn yank_current_code(&self) -> Option<String> {
        let filtered = self.filtered_entries();
        filtered.get(self.selected_index)
            .and_then(|entry| crate::otp::generate_code(entry).ok())
            .map(|code| code.value)
    }

    fn filter_by_property<'a>(&'a self, prop: &str, query: &str) -> Vec<&'a Entry> {
        let mut matcher = Matcher::new(Config::DEFAULT);
        self.entries.iter()
            .filter(|e| {
                if prop.starts_with("iss") {
                    return nucleo_match(query, &e.issuer, &mut matcher);
                }
                if prop.starts_with("nam") {
                    return nucleo_match(query, &e.name, &mut matcher);
                }
                if prop.starts_with("not") {
                    return nucleo_match(query, &e.note, &mut matcher);
                }
                if prop.starts_with("fav") {
                    let fav_str = if e.favorite { "true" } else { "false" };
                    return nucleo_match(query, fav_str, &mut matcher);
                }
                if prop.starts_with("typ") {
                    return nucleo_match(query, &e.entry_type, &mut matcher);
                }
                false
            })
            .collect()
    }

    fn fuzzy_filter_entries<'a, F>(&'a self, query: &str, field_extractor: F) -> Vec<&'a Entry>
    where
        F: Fn(&Entry) -> &str,
    {
        let mut matcher = Matcher::new(Config::DEFAULT);
        let pattern = Pattern::new(query, CaseMatching::Ignore, Normalization::Smart, nucleo::pattern::AtomKind::Fuzzy);
        let mut buf = Vec::new();
        self.entries.iter()
            .filter(|e| {
                buf.clear();
                let text = Utf32Str::new(field_extractor(e), &mut buf);
                pattern.indices(text, &mut matcher, &mut Vec::new()).is_some()
            })
            .collect()
    }
}

fn nucleo_match(pattern: &str, text: &str, matcher: &mut Matcher) -> bool {
    let p = Pattern::new(pattern, CaseMatching::Ignore, Normalization::Smart, nucleo::pattern::AtomKind::Fuzzy);
    let mut buf = Vec::new();
    let text = Utf32Str::new(text, &mut buf);
    p.indices(text, matcher, &mut Vec::new()).is_some()
}

impl Default for App {
    fn default() -> Self {
        Self::new()
    }
}

use std::io;
