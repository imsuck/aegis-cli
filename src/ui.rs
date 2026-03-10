use ratatui::{
    Frame,
    layout::{Constraint, Direction, Layout, Rect},
    style::{Color, Style},
    widgets::{Block, Borders, Clear, List, ListItem, Paragraph},
};
use crate::app::App;
use crate::otp::generate_code;

impl App {
    pub fn draw(&self, frame: &mut Frame) {
        let chunks = Layout::default()
            .direction(Direction::Vertical)
            .constraints([
                Constraint::Length(3),
                Constraint::Min(0),
                Constraint::Length(3),
            ])
            .split(frame.area());
        
        // Title
        let title = Paragraph::new("Aegis CLI - 2FA Vault")
            .style(Style::default().fg(Color::Cyan))
            .block(Block::default().borders(Borders::ALL));
        frame.render_widget(title, chunks[0]);
        
        // Entry list
        Self::render_entry_list(frame, self, chunks[1]);
        
        // Status bar
        Self::render_status_bar(frame, self, chunks[2]);
        
        // Search popup if in search mode
        if self.search_mode {
            Self::render_search_popup(frame, self);
        }
    }

    fn render_entry_list(frame: &mut Frame, app: &App, area: Rect) {
        let filtered = app.filtered_entries();
        
        let items: Vec<ListItem> = filtered.iter().enumerate().map(|(i, entry)| {
            let code = if app.show_code && i == app.selected_index {
                if let Ok(code) = generate_code(entry) {
                    // Show asterisks instead of actual code for security
                    let masked = "*".repeat(code.value.len());
                    format!(" [{} | {}s]", masked, code.period_remaining)
                } else {
                    String::new()
                }
            } else {
                String::new()
            };
            
            let note_preview = if entry.note.len() > 30 {
                format!("{}...", &entry.note[..30])
            } else {
                entry.note.clone()
            };
            
            let content = format!(
                "{} | {} | {}{}",
                entry.issuer,
                entry.name,
                if note_preview.is_empty() { "(no note)" } else { &note_preview },
                code
            );
            
            ListItem::new(content)
        }).collect();
        
        let list = List::new(items)
            .block(Block::default().title("Entries").borders(Borders::ALL));
        
        frame.render_widget(list, area);
    }

    fn render_status_bar(frame: &mut Frame, app: &App, area: Rect) {
        let status = if app.search_mode {
            format!("Search: {} (Esc to cancel)", app.search_query)
        } else if app.show_code {
            "Press 'y' to yank code, 'c' to hide, 'q' to quit".to_string()
        } else {
            "j/k: navigate | /: search | c: show code | q: quit".to_string()
        };
        
        let status_widget = Paragraph::new(status)
            .style(Style::default().fg(Color::White));
        
        frame.render_widget(status_widget, area);
    }

    fn render_search_popup(frame: &mut Frame, app: &App) {
        let area = Self::centered_rect(60, 20, frame.area());
        frame.render_widget(Clear, area);
        
        let help_text = "Properties: %issuer, %name, %note, %favorite, %type";
        let input = Paragraph::new(app.search_query.as_str())
            .block(Block::default()
                .title("Search (/)")
                .title_bottom(help_text)
                .borders(Borders::ALL))
            .style(Style::default().fg(Color::Yellow));
        
        frame.render_widget(input, area);
        frame.set_cursor_position((area.x + 1, area.y + 1));
    }

    fn centered_rect(percent_x: u16, percent_y: u16, area: Rect) -> Rect {
        let popup_layout = Layout::default()
            .direction(Direction::Vertical)
            .constraints([
                Constraint::Percentage((100 - percent_y) / 2),
                Constraint::Percentage(percent_y),
                Constraint::Percentage((100 - percent_y) / 2),
            ])
            .split(area);
        
        Layout::default()
            .direction(Direction::Horizontal)
            .constraints([
                Constraint::Percentage((100 - percent_x) / 2),
                Constraint::Percentage(percent_x),
                Constraint::Percentage((100 - percent_x) / 2),
            ])
            .split(popup_layout[1])[1]
    }
}
