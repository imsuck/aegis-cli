mod vault;
mod decrypt;
mod otp;
mod app;
mod ui;

use std::io::{self, stdout, Write};
use crossterm::{
    event::{self, DisableMouseCapture, EnableMouseCapture, Event, KeyCode},
    execute,
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
};
use ratatui::prelude::*;
use app::App;
use ui::render;
use decrypt::decrypt_vault;
use std::fs;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Load vault file first
    let vault_path = std::env::args().nth(1)
        .unwrap_or_else(|| "vault.json".to_string());
    let vault_content = fs::read_to_string(&vault_path)?;

    // Prompt for password BEFORE setting up TUI
    print!("Password: ");
    io::stdout().flush()?;
    let mut password = String::new();
    io::stdin().read_line(&mut password)?;
    let password = password.trim().to_string();

    // Decrypt vault
    let decrypted = decrypt_vault(&vault_content, &password)?;

    // Now setup terminal for TUI
    enable_raw_mode()?;
    let mut stdout = stdout();
    execute!(stdout, EnterAlternateScreen, EnableMouseCapture)?;
    let backend = CrosstermBackend::new(stdout);
    let mut terminal = Terminal::new(backend)?;

    // Create app
    let mut app = App::new();
    app.set_entries(decrypted.entries);
    app.set_password(password);

    // Run TUI
    let result = run_app(&mut terminal, &mut app);

    // Restore terminal
    disable_raw_mode()?;
    execute!(
        terminal.backend_mut(),
        LeaveAlternateScreen,
        DisableMouseCapture
    )?;
    terminal.show_cursor()?;

    if let Err(err) = result {
        eprintln!("Error: {:?}", err);
    }

    Ok(())
}

fn run_app<B: Backend>(terminal: &mut Terminal<B>, app: &mut App) -> io::Result<()> {
    use std::time::Duration;
    
    loop {
        terminal.draw(|f| render(f, app))?;

        // Poll for events with 100ms timeout for smooth countdown refresh
        if event::poll(Duration::from_millis(100))? {
            if let Event::Key(key) = event::read()? {
                match key.code {
                    KeyCode::Char('q') => app.running = false,
                    KeyCode::Char('j') | KeyCode::Down => {
                        if app.selected_index < app.filtered_entries().len().saturating_sub(1) {
                            app.selected_index += 1;
                        }
                    }
                    KeyCode::Char('k') | KeyCode::Up => {
                        if app.selected_index > 0 {
                            app.selected_index -= 1;
                        }
                    }
                    KeyCode::Char('/') => {
                        app.search_mode = true;
                        app.search_query.clear();
                    }
                    KeyCode::Esc if app.search_mode => {
                        app.search_mode = false;
                        app.search_query.clear();
                    }
                    KeyCode::Enter if app.search_mode => {
                        app.search_mode = false;
                    }
                    KeyCode::Char(c) if app.search_mode => {
                        app.search_query.push(c);
                    }
                    KeyCode::Char('c') => {
                        app.show_code = !app.show_code;
                    }
                    KeyCode::Char('y') if app.show_code => {
                        if let Some(code) = app.yank_current_code() {
                            if let Ok(mut clipboard) = arboard::Clipboard::new() {
                                let _ = clipboard.set_text(&code);
                            }
                        }
                    }
                    _ => {}
                }
            }
        }

        if !app.running {
            break;
        }
    }
    
    Ok(())
}
