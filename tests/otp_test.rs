use aegis_cli::otp::generate_code;
use aegis_cli::vault::{Entry, EntryInfo, TotpInfo};
use uuid::Uuid;

#[test]
fn test_generate_totp_code() {
    let entry = Entry {
        entry_type: "totp".to_string(),
        uuid: Uuid::new_v4(),
        name: "Test".to_string(),
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
    
    let code = generate_code(&entry);
    assert!(code.is_ok());
    let code = code.unwrap();
    assert_eq!(code.value.len(), 6);
    assert!(code.period_remaining > 0);
    assert!(code.period_remaining <= 30);
}
