package io.github.aadegtyarev.sitewrap

import android.content.Context
import android.util.Base64
import java.net.URI
import java.security.SecureRandom
import javax.crypto.SecretKeyFactory
import javax.crypto.spec.PBEKeySpec

data class Config(val url: String, val zoom: Int, val mode: String, val language: String) {
    fun error(): String? {
        val uri = runCatching { URI(url) }.getOrNull()
        if (uri == null || uri.scheme !in listOf("http", "https") || uri.host.isNullOrBlank() || uri.userInfo != null) return "invalidUrl"
        if (zoom !in 30..500) return "invalidZoom"
        if (mode !in listOf("window", "kiosk")) return "invalidMode"
        if (language !in listOf("en", "de", "kk", "it")) return "invalidLanguage"
        return null
    }
}

class ConfigStore(context: Context) {
    private val prefs = context.getSharedPreferences("sitewrap", Context.MODE_PRIVATE)
    fun load(): Config? {
        if (!prefs.contains("url")) return null
        if (!prefs.contains("salt") || !prefs.contains("hash")) return null
        return Config(prefs.getString("url", "") ?: "", prefs.getInt("zoom", 100), prefs.getString("mode", "window") ?: "window", prefs.getString("language", "en") ?: "en")
            .takeIf { it.error() == null }
    }
    fun hasConfig() = prefs.contains("url")
    fun verify(password: String): Boolean {
        val salt = prefs.getString("salt", null) ?: error("Missing password salt")
        val expected = prefs.getString("hash", null) ?: error("Missing password verifier")
        val actual = derive(password, Base64.decode(salt, Base64.NO_WRAP))
        return java.security.MessageDigest.isEqual(actual, Base64.decode(expected, Base64.NO_WRAP))
    }
    fun save(config: Config, password: String?): Boolean {
        require(config.error() == null)
        val edit = prefs.edit().putString("url", config.url).putInt("zoom", config.zoom)
            .putString("mode", config.mode).putString("language", config.language)
        if (password != null) {
            require(password.length >= 12)
            val salt = ByteArray(16).also(SecureRandom()::nextBytes)
            edit.putString("salt", Base64.encodeToString(salt, Base64.NO_WRAP))
                .putString("hash", Base64.encodeToString(derive(password, salt), Base64.NO_WRAP))
        }
        return edit.commit()
    }
    private fun derive(password: String, salt: ByteArray): ByteArray {
        val spec = PBEKeySpec(password.toCharArray(), salt, 210_000, 256)
        return try { SecretKeyFactory.getInstance("PBKDF2WithHmacSHA256").generateSecret(spec).encoded }
        finally { spec.clearPassword() }
    }
}
