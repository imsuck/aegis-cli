use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Vault {
    pub version: u32,
    pub header: Header,
    pub db: String,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Header {
    pub slots: Vec<Slot>,
    pub params: Params,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Slot {
    #[serde(rename = "type")]
    pub slot_type: u8,
    pub uuid: Uuid,
    pub key: String,
    pub key_params: KeyParams,
    #[serde(default)]
    pub n: u32,
    #[serde(default)]
    pub r: u32,
    #[serde(default)]
    pub p: u32,
    #[serde(default)]
    pub salt: String,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Params {
    pub nonce: String,
    pub tag: String,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct KeyParams {
    pub nonce: String,
    pub tag: String,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct VaultContent {
    pub version: u32,
    pub entries: Vec<Entry>,
    #[serde(default)]
    pub groups: Vec<Group>,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Entry {
    #[serde(rename = "type")]
    pub entry_type: String,
    pub uuid: Uuid,
    pub name: String,
    pub issuer: String,
    #[serde(default)]
    pub note: String,
    #[serde(default)]
    pub favorite: bool,
    pub info: EntryInfo,
    #[serde(default)]
    pub groups: Vec<Uuid>,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
#[serde(untagged)]
pub enum EntryInfo {
    Totp(TotpInfo),
    Hotp(HotpInfo),
    Steam(SteamInfo),
    Motp(MotpInfo),
    Yandex(YandexInfo),
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct TotpInfo {
    pub secret: String,
    pub algo: String,
    pub digits: u32,
    pub period: u32,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct HotpInfo {
    pub secret: String,
    pub algo: String,
    pub digits: u32,
    pub counter: u64,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct SteamInfo {
    pub secret: String,
    pub digits: u32,
    pub period: u32,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct MotpInfo {
    pub secret: String,
    pub digits: u32,
    pub period: u32,
    pub pin: String,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct YandexInfo {
    pub secret: String,
    pub digits: u32,
    pub period: u32,
    pub pin: String,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Group {
    pub uuid: Uuid,
    pub name: String,
}
