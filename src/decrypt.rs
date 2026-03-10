use crate::vault::{Vault, VaultContent, Slot};
use aes_gcm::{Aes256Gcm, KeyInit, Nonce, aead::Aead};
use scrypt::{scrypt, Params as ScryptParams};
use base64::{Engine as _, engine::general_purpose::STANDARD as BASE64};

pub fn decrypt_vault(vault_json: &str, password: &str) -> Result<VaultContent, DecryptError> {
    let vault: Vault = serde_json::from_str(vault_json)
        .map_err(|e| DecryptError::ParseError(e.to_string()))?;
    
    // Find password slots (type 1)
    let password_slots: Vec<&Slot> = vault.header.slots
        .iter()
        .filter(|s| s.slot_type == 1)
        .collect();
    
    if password_slots.is_empty() {
        return Err(DecryptError::NoPasswordSlot);
    }
    
    // Try each password slot
    let master_key = derive_master_key(&password_slots, password.as_bytes())?;
    
    // Decrypt the vault contents
    let encrypted_content = BASE64.decode(&vault.db)
        .map_err(|e| DecryptError::Base64Error(e.to_string()))?;
    
    let nonce = hex::decode(&vault.header.params.nonce)
        .map_err(|e| DecryptError::HexError(e.to_string()))?;
    let tag = hex::decode(&vault.header.params.tag)
        .map_err(|e| DecryptError::HexError(e.to_string()))?;
    
    let cipher = Aes256Gcm::new_from_slice(&master_key)
        .map_err(|e| DecryptError::CipherError(e.to_string()))?;
    
    let mut decrypted_data = encrypted_content;
    decrypted_data.extend_from_slice(&tag);
    
    let nonce = Nonce::from_slice(&nonce);
    let decrypted = cipher.decrypt(nonce, decrypted_data.as_ref())
        .map_err(|e| DecryptError::DecryptionError(e.to_string()))?;
    
    let content: VaultContent = serde_json::from_slice(&decrypted)
        .map_err(|e| DecryptError::ParseError(e.to_string()))?;
    
    Ok(content)
}

fn derive_master_key(slots: &[&Slot], password: &[u8]) -> Result<Vec<u8>, DecryptError> {
    for slot in slots {
        let salt = hex::decode(&slot.salt)
            .map_err(|e| DecryptError::HexError(e.to_string()))?;
        
        // scrypt uses log_n (log2 of N), so N=32768 -> log_n=15
        let log_n = (slot.n as f64).log2() as u8;
        let scrypt_params = ScryptParams::new(log_n, slot.r, slot.p, 32)
            .map_err(|e| DecryptError::ScryptError(e.to_string()))?;
        
        let mut derived_key = [0u8; 32];
        scrypt(password, &salt, &scrypt_params, &mut derived_key)
            .map_err(|e| DecryptError::ScryptError(e.to_string()))?;
        
        // Try to decrypt the master key with this derived key
        let key_bytes = hex::decode(&slot.key)
            .map_err(|e| DecryptError::HexError(e.to_string()))?;
        let nonce = hex::decode(&slot.key_params.nonce)
            .map_err(|e| DecryptError::HexError(e.to_string()))?;
        let tag = hex::decode(&slot.key_params.tag)
            .map_err(|e| DecryptError::HexError(e.to_string()))?;
        
        let cipher = Aes256Gcm::new_from_slice(&derived_key)
            .map_err(|e| DecryptError::CipherError(e.to_string()))?;
        
        let mut encrypted_key = key_bytes;
        encrypted_key.extend_from_slice(&tag);
        
        let nonce = Nonce::from_slice(&nonce);
        if let Ok(master_key) = cipher.decrypt(nonce, encrypted_key.as_ref()) {
            return Ok(master_key);
        }
    }
    
    Err(DecryptError::InvalidPassword)
}

#[derive(Debug)]
pub enum DecryptError {
    ParseError(String),
    NoPasswordSlot,
    Base64Error(String),
    HexError(String),
    CipherError(String),
    DecryptionError(String),
    ScryptError(String),
    InvalidPassword,
}

impl std::fmt::Display for DecryptError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{:?}", self)
    }
}

impl std::error::Error for DecryptError {}
