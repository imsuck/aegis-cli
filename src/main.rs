mod vault;
mod decrypt;
mod otp;
mod app;
mod ui;

use std::io::{self, stdout, Write};
use crossterm::{
    event::{self, Event, KeyCode, KeyEventKind},
    execute,
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
};
use ratatui::{backend::CrosstermBackend, Terminal};
use app::App;
use decrypt::decrypt_vault;
use std::fs;

/// Prompt for password with redacted input (shows asterisks)
fn prompt_password() -> io::Result<String> {
    print!("Password: ");
    stdout().flush()?;
    
    enable_raw_mode()?;
    let mut password = String::new();
    
    loop {
        if let Event::Key(key) = event::read()? {
            if key.kind != KeyEventKind::Press {
                continue;
            }
            
            match key.code {
                KeyCode::Enter => {
                    break;
                }
                KeyCode::Char(c) if !key.modifiers.contains(crossterm::event::KeyModifiers::CONTROL) => {
                    password.push(c);
                    print!("*");
                    stdout().flush()?;
                }
                KeyCode::Backspace => {
                    if !password.is_empty() {
                        password.pop();
                        // Move cursor back, erase char, move back again
                        print!("\x08 \x08");
                        stdout().flush()?;
                    }
                }
                KeyCode::Esc => {
                    // Esc to cancel
                    println!("\nCancelled.");
                    disable_raw_mode()?;
                    return Err(io::Error::new(io::ErrorKind::Interrupted, "Password input cancelled"));
                }
                _ => {}
            }
        }
    }
    
    disable_raw_mode()?;
    println!();
    
    Ok(password)
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Load vault file first
    let vault_path = std::env::args().nth(1)
        .unwrap_or_else(|| "vault.json".to_string());
    let vault_content = fs::read_to_string(&vault_path)?;

    // Prompt for password with redacted input
    let password = prompt_password()?;

    // Decrypt vault
    let decrypted = decrypt_vault(&vault_content, &password)?;

    // Create app
    let mut app = App::new();
    app.set_entries(decrypted.entries);
    app.set_password(password);

    // Setup terminal
    enable_raw_mode()?;
    let mut stdout = stdout();
    execute!(stdout, EnterAlternateScreen)?;
    let backend = CrosstermBackend::new(stdout);
    let mut terminal = Terminal::new(backend)?;

    // Run the application
    let res = app.run(&mut terminal);

    // Restore terminal
    disable_raw_mode()?;
    execute!(terminal.backend_mut(), LeaveAlternateScreen)?;

    // Check for errors from the app
    if let Err(err) = res {
        eprintln!("Error: {:?}", err);
    }

    Ok(())
}
