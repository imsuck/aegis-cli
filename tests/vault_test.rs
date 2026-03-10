use aegis_cli::vault::Vault;

#[test]
fn test_parse_encrypted_vault_structure() {
    let json = r#"{
        "version": 1,
        "header": {
            "slots": [{
                "type": 1,
                "uuid": "01234567-89ab-cdef-0123-456789abcdef",
                "key": "abcdef0123456789",
                "key_params": {
                    "nonce": "0123456789abcdef01234567",
                    "tag": "0123456789abcdef0123456789abcdef"
                },
                "n": 32768,
                "r": 8,
                "p": 1,
                "salt": "0123456789abcdef"
            }],
            "params": {
                "nonce": "0123456789abcdef01234567",
                "tag": "0123456789abcdef0123456789abcdef"
            }
        },
        "db": "base64encodedcontent"
    }"#;
    
    let vault: Vault = serde_json::from_str(json).unwrap();
    assert_eq!(vault.version, 1);
    assert_eq!(vault.header.slots.len(), 1);
}
