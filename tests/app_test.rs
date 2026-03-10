use aegis_cli::app::App;

#[test]
fn test_app_creation() {
    let app = App::new();
    // New app should have no entries and selection at 0
    assert!(app.entries.is_empty());
    assert_eq!(app.get_selected_index(), 0);
}

#[test]
fn test_selection_persists_across_search() {
    use aegis_cli::vault::{Entry, EntryInfo, TotpInfo};
    use uuid::Uuid;
    
    let mut app = App::new();
    
    // Add two entries
    let entry1 = Entry {
        entry_type: "totp".to_string(),
        uuid: Uuid::new_v4(),
        name: "Alice".to_string(),
        issuer: "Google".to_string(),
        note: String::new(),
        favorite: false,
        info: EntryInfo::Totp(TotpInfo {
            secret: "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ".to_string(),
            algo: "SHA1".to_string(),
            digits: 6,
            period: 30,
        }),
        groups: vec![],
    };
    let entry2 = Entry {
        entry_type: "totp".to_string(),
        uuid: Uuid::new_v4(),
        name: "Bob".to_string(),
        issuer: "GitHub".to_string(),
        note: String::new(),
        favorite: false,
        info: EntryInfo::Totp(TotpInfo {
            secret: "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ".to_string(),
            algo: "SHA1".to_string(),
            digits: 6,
            period: 30,
        }),
        groups: vec![],
    };
    
    let uuid1 = entry1.uuid;
    app.set_entries(vec![entry1, entry2]);
    
    // Move to second entry
    app.move_selection(1);
    assert_eq!(app.get_selected_index(), 1);
    
    // Search should filter but keep selection on same entry by UUID
    app.search_query = "github".to_string();
    let filtered = app.filtered_entries();
    assert_eq!(filtered.len(), 1);
    // Selection should now show index 0 (only one result) but still be Bob
    assert_eq!(app.get_selected_index(), 0);
    assert_eq!(app.get_selected_entry().unwrap().name, "Bob");
    
    // Clear search - selection should return to original position
    app.search_query.clear();
    assert_eq!(app.get_selected_index(), 1);
    assert_eq!(app.get_selected_entry().unwrap().name, "Bob");
}
