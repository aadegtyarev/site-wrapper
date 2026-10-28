package io.github.aadegtyarev.sitewrap

import org.junit.Assert.assertEquals
import org.junit.Assert.assertNull
import org.junit.Test

class ConfigTest {
    @Test fun validDefaults() {
        assertNull(Config("https://example.org/path", 100, "window", "en").error())
    }
    @Test fun invalidInputs() {
        assertEquals("invalidUrl", Config("https://user:pass@example.org", 100, "window", "en").error())
        assertEquals("invalidUrl", Config("file:///tmp/x", 100, "window", "en").error())
        assertEquals("invalidZoom", Config("https://example.org", 29, "window", "en").error())
        assertEquals("invalidMode", Config("https://example.org", 100, "other", "en").error())
        assertEquals("invalidLanguage", Config("https://example.org", 100, "window", "fr").error())
    }
    @Test fun translationsCoverEnglishKeys() {
        for (language in listOf("en", "de", "kk", "it")) {
            for (key in listOf("settings", "url", "zoom", "password", "save", "invalidUrl", "loadFailed", "unsupportedLink")) {
                assertEquals(false, Words.get(language, key) == key)
            }
        }
    }
}
