use aegis_cli::app::App;

#[test]
fn test_app_creation() {
    let app = App::new();
    assert_eq!(app.selected_index, 0);
    assert!(app.entries.is_empty());
}
