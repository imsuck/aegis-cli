use aegis_cli::decrypt::decrypt_vault;

#[test]
fn test_decrypt_test_vault() {
    let vault_path = "test/resources/aegis_encrypted.json";
    let content = std::fs::read_to_string(vault_path).unwrap();
    
    // Password is "test" per the decrypt.py example
    let result = decrypt_vault(&content, "test");
    assert!(result.is_ok());
    
    let vault_content = result.unwrap();
    assert!(!vault_content.entries.is_empty());
}
