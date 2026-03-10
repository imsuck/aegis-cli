use crate::vault::Entry;
use zeroize::Zeroizing;

#[derive(Debug)]
pub struct App {
    pub entries: Vec<Entry>,
    pub selected_index: usize,
    pub password: Zeroizing<String>,
    pub search_query: String,
    pub search_mode: bool,
    pub show_code: bool,
    pub running: bool,
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
            running: true,
        }
    }
    
    pub fn set_entries(&mut self, entries: Vec<Entry>) {
        self.entries = entries;
    }
    
    pub fn set_password(&mut self, password: String) {
        self.password = Zeroizing::new(password);
    }
    
    pub fn filtered_entries(&self) -> Vec<&Entry> {
        if self.search_query.is_empty() {
            return self.entries.iter().collect();
        }
        
        let query = self.search_query.to_lowercase();
        self.entries.iter()
            .filter(|e| {
                e.issuer.to_lowercase().contains(&query) ||
                e.name.to_lowercase().contains(&query) ||
                e.note.to_lowercase().contains(&query)
            })
            .collect()
    }
}

impl Default for App {
    fn default() -> Self {
        Self::new()
    }
}
