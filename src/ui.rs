use crate::app::App;
use crate::otp::generate_code;
use ratatui::{
    layout::{Constraint, Direction, Layout, Rect},
    style::{Color, Modifier, Style},
    widgets::{Block, Borders, Clear, Paragraph, Row, Table},
    Frame,
};

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

        // Entry table
        self.render_entry_table(frame, chunks[1]);

        // Status bar
        Self::render_status_bar(frame, self, chunks[2]);

        // Search popup if in search mode
        if self.search_mode {
            Self::render_search_popup(frame, self);
        }
    }

    fn render_entry_table(&self, frame: &mut Frame, area: Rect) {
        let filtered = self.filtered_entries();
        let selected_index = self.get_selected_index();

        // Pre-compute codes for all entries in a single pass
        let codes: Vec<Option<(String, u64)>> = filtered
            .iter()
            .map(|entry| {
                generate_code(entry)
                    .ok()
                    .map(|code| (code.value, code.period_remaining))
            })
            .collect();

        // Calculate max widths for proper column alignment (single pass)
        let (max_issuer_len, max_name_len) =
            filtered.iter().fold((0, 0), |(max_iss, max_name), e| {
                (
                    max_iss.max(e.issuer.len().min(30)),
                    max_name.max(e.name.len().min(20)),
                )
            });

        let rows: Vec<Row> = filtered
            .iter()
            .enumerate()
            .map(|(i, entry)| {
                let code_info = if self.show_code && i == selected_index {
                    if let Some((value, period)) = &codes[i] {
                        let masked = "*".repeat(value.len());
                        format!("{} | {}s", masked, period)
                    } else {
                        "****** | --s".to_owned()
                    }
                } else {
                    String::new()
                };

                let note_preview = if entry.note.len() > 20 {
                    format!("{}...", &entry.note[..20])
                } else {
                    entry.note.clone()
                };

                let style = if i == selected_index {
                    Style::default()
                        .fg(Color::Yellow)
                        .add_modifier(Modifier::BOLD)
                } else {
                    Style::default()
                };

                Row::new(vec![
                    entry.issuer.clone(),
                    entry.name.clone(),
                    note_preview,
                    code_info,
                ])
                .style(style)
            })
            .collect();

        // Create table with column widths
        let widths = vec![
            Constraint::Length(max_issuer_len as u16 + 2),
            Constraint::Length(max_name_len as u16 + 2),
            Constraint::Length(22), // note column
            Constraint::Length(15), // code column
        ];

        let table = Table::new(rows, widths)
            .header(
                Row::new(vec!["Issuer", "Name", "Note", "Code"])
                    .style(
                        Style::default()
                            .fg(Color::White)
                            .add_modifier(Modifier::BOLD),
                    )
                    .bottom_margin(1),
            )
            .block(Block::default().title("Entries").borders(Borders::ALL))
            .row_highlight_style(Style::default().add_modifier(Modifier::REVERSED))
            .highlight_symbol(">> ");

        frame.render_widget(table, area);
    }

    fn render_status_bar(frame: &mut Frame, app: &App, area: Rect) {
        let status = if app.search_mode {
            format!("Search: {} (Esc to cancel)", app.search_query)
        } else if app.show_code {
            "y: yank code | c: hide | q: quit".to_owned()
        } else {
            "j/k: navigate | /: search | c: show code | y: yank | q: quit".to_owned()
        };

        let status_widget = Paragraph::new(status).style(Style::default().fg(Color::White));

        frame.render_widget(status_widget, area);
    }

    fn render_search_popup(frame: &mut Frame, app: &App) {
        let area = Self::centered_rect(60, 20, frame.area());
        frame.render_widget(Clear, area);

        let help_text = "Properties: %issuer, %name, %note, %favorite, %type";
        let input = Paragraph::new(app.search_query.as_str())
            .block(
                Block::default()
                    .title("Search (/)")
                    .title_bottom(help_text)
                    .borders(Borders::ALL),
            )
            .style(Style::default().fg(Color::Yellow));

        frame.render_widget(input, area);
        // Position cursor at the end of the current query text
        let cursor_x = area.x + 1 + app.search_query.len() as u16;
        let cursor_y = area.y + 1;
        frame.set_cursor_position((cursor_x, cursor_y));
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
