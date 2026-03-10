use crate::vault::{Entry, EntryInfo, TotpInfo};
use totp_rs::{Algorithm, Secret, TOTP};
use chrono::Utc;

#[derive(Debug)]
pub struct Code {
    pub value: String,
    pub period_remaining: u64,
}

pub fn generate_code(entry: &Entry) -> Result<Code, OtpError> {
    match &entry.info {
        EntryInfo::Totp(info) => generate_totp_code(entry, info),
        EntryInfo::Hotp(info) => generate_hotp_code(entry, info),
        EntryInfo::Steam(info) => generate_steam_code(entry, info),
        EntryInfo::Motp(info) => generate_motp_code(entry, info),
        EntryInfo::Yandex(info) => generate_yandex_code(entry, info),
    }
}

fn generate_totp_code(_entry: &Entry, info: &TotpInfo) -> Result<Code, OtpError> {
    let secret = Secret::Encoded(info.secret.clone())
        .to_raw()
        .map_err(|e| OtpError::InvalidSecret(e.to_string()))?;
    
    let algo = parse_algorithm(&info.algo)?;
    
    let totp = TOTP::new(
        algo,
        info.digits as usize,
        0,
        info.period as u64,
        secret.to_bytes().map_err(|e| OtpError::InvalidSecret(e.to_string()))?,
    ).map_err(|e| OtpError::TotpError(e.to_string()))?;
    
    let timestamp = Utc::now().timestamp() as u64;
    let value = totp.generate(timestamp);
    
    let period_remaining = (info.period as u64) - (timestamp % (info.period as u64));
    
    Ok(Code {
        value,
        period_remaining,
    })
}

fn generate_hotp_code(_entry: &Entry, _info: &crate::vault::HotpInfo) -> Result<Code, OtpError> {
    Err(OtpError::NotImplemented("HOTP counter tracking not implemented".to_string()))
}

fn generate_steam_code(_entry: &Entry, _info: &crate::vault::SteamInfo) -> Result<Code, OtpError> {
    Err(OtpError::NotImplemented("Steam OTP not fully implemented".to_string()))
}

fn generate_motp_code(_entry: &Entry, _info: &crate::vault::MotpInfo) -> Result<Code, OtpError> {
    Err(OtpError::NotImplemented("MOTP not fully implemented".to_string()))
}

fn generate_yandex_code(_entry: &Entry, _info: &crate::vault::YandexInfo) -> Result<Code, OtpError> {
    Err(OtpError::NotImplemented("Yandex OTP not fully implemented".to_string()))
}

fn parse_algorithm(algo: &str) -> Result<Algorithm, OtpError> {
    match algo.to_uppercase().as_str() {
        "SHA1" | "SHA-1" => Ok(Algorithm::SHA1),
        "SHA256" | "SHA-256" => Ok(Algorithm::SHA256),
        "SHA512" | "SHA-512" => Ok(Algorithm::SHA512),
        _ => Err(OtpError::InvalidAlgorithm(algo.to_string())),
    }
}

#[derive(Debug)]
pub enum OtpError {
    InvalidSecret(String),
    InvalidAlgorithm(String),
    TotpError(String),
    NotImplemented(String),
}

impl std::fmt::Display for OtpError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{:?}", self)
    }
}

impl std::error::Error for OtpError {}
