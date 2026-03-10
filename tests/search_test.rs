use aegis_cli::app::App;
use aegis_cli::vault::{Entry, EntryInfo, TotpInfo};
use uuid::Uuid;

#[test]
fn test_property_search() {
    let mut app = App::new();
    
    let entry1 = Entry {
        entry_type: "totp".to_string(),
        uuid: Uuid::new_v4(),
        name: "Alice".to_string(),
        issuer: "Google".to_string(),
        note: "Work account".to_string(),
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
        note: "Personal".to_string(),
        favorite: true,
        info: EntryInfo::Totp(TotpInfo {
            secret: "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ".to_string(),
            algo: "SHA1".to_string(),
            digits: 6,
            period: 30,
        }),
        groups: vec![],
    };
    
    app.set_entries(vec![entry1.clone(), entry2.clone()]);
    
    // Default search (issuer)
    app.search_query = "goo".to_string();
    let filtered = app.filtered_entries();
    assert_eq!(filtered.len(), 1);
    assert_eq!(filtered[0].issuer, "Google");
    
    // Property prefix search
    app.search_query = "%name ali".to_string();
    let filtered = app.filtered_entries();
    assert_eq!(filtered.len(), 1);
    assert_eq!(filtered[0].name, "Alice");
    
    // Property prefix search for note
    app.search_query = "%note work".to_string();
    let filtered = app.filtered_entries();
    assert_eq!(filtered.len(), 1);
}
