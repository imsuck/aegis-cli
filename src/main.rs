mod vault;
mod decrypt;
mod otp;
mod app;
mod ui;

use std::io::{self, stdout, Write};
use crossterm::{
    execute,
    terminal::{disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen},
};
use ratatui::{backend::CrosstermBackend, Terminal};
use app::App;
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
