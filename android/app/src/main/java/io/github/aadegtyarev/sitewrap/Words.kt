package io.github.aadegtyarev.sitewrap

object Words {
    private val en = mapOf(
        "settings" to "Settings", "unlock" to "Unlock settings", "password" to "Settings password", "open" to "Unlock",
        "url" to "Website URL", "zoom" to "Page zoom (%)", "mode" to "Display mode", "language" to "Language",
        "newPassword" to "New password", "confirm" to "Confirm password", "save" to "Save and open", "cancel" to "Cancel",
        "window" to "Window", "kiosk" to "Kiosk", "version" to "Version", "retry" to "Retry",
        "invalidUrl" to "Enter a full HTTP(S) URL without credentials.", "invalidZoom" to "Choose a zoom between 30% and 500%.",
        "invalidMode" to "Choose a display mode.", "invalidLanguage" to "Choose a language.",
        "passwordShort" to "Use at least 12 characters.", "passwordMismatch" to "Passwords do not match.",
        "incorrectPassword" to "Incorrect password. Try again.", "saveFailed" to "Changes were not saved. Check device storage and try again.",
        "loadFailed" to "The site could not load. Check your connection or change the URL in settings.",
        "configFailed" to "Settings could not be read. Check device storage or reinstall Sitewrap.",
        "extensionFailed" to "Page zoom could not be enabled. Restart or reinstall Sitewrap.",
        "unsupportedLink" to "This link needs another app or a web address."
    )
    private val de = mapOf(
        "settings" to "Einstellungen", "unlock" to "Einstellungen entsperren", "password" to "Passwort für Einstellungen", "open" to "Entsperren",
        "url" to "Website-Adresse", "zoom" to "Seitenzoom (%)", "mode" to "Anzeigemodus", "language" to "Sprache",
        "newPassword" to "Neues Passwort", "confirm" to "Passwort bestätigen", "save" to "Speichern und öffnen", "cancel" to "Abbrechen",
        "window" to "Fenster", "kiosk" to "Kiosk", "version" to "Version", "retry" to "Erneut versuchen",
        "invalidUrl" to "Vollständige HTTP(S)-Adresse ohne Zugangsdaten eingeben.", "invalidZoom" to "Zoom zwischen 30 % und 500 % wählen.",
        "invalidMode" to "Anzeigemodus wählen.", "invalidLanguage" to "Sprache wählen.",
        "passwordShort" to "Mindestens 12 Zeichen verwenden.", "passwordMismatch" to "Die Passwörter stimmen nicht überein.",
        "incorrectPassword" to "Falsches Passwort. Erneut versuchen.", "saveFailed" to "Änderungen nicht gespeichert. Gerätespeicher prüfen und erneut versuchen.",
        "loadFailed" to "Website konnte nicht geladen werden. Verbindung prüfen oder Adresse in den Einstellungen ändern.",
        "configFailed" to "Einstellungen konnten nicht gelesen werden. Gerätespeicher prüfen oder Sitewrap neu installieren.",
        "extensionFailed" to "Seitenzoom konnte nicht aktiviert werden. Sitewrap neu starten oder neu installieren.",
        "unsupportedLink" to "Für diesen Link wird eine andere App oder eine Webadresse benötigt."
    )
    private val kk = mapOf(
        "settings" to "Баптаулар", "unlock" to "Баптауларды ашу", "password" to "Баптаулар құпиясөзі", "open" to "Ашу",
        "url" to "Сайт мекенжайы", "zoom" to "Бет масштабы (%)", "mode" to "Көрсету режимі", "language" to "Тіл",
        "newPassword" to "Жаңа құпиясөз", "confirm" to "Құпиясөзді растау", "save" to "Сақтау және ашу", "cancel" to "Бас тарту",
        "window" to "Терезе", "kiosk" to "Киоск", "version" to "Нұсқа", "retry" to "Қайталау",
        "invalidUrl" to "Деректемелерсіз толық HTTP(S) мекенжайын енгізіңіз.", "invalidZoom" to "Масштабты 30%-дан 500%-ға дейін таңдаңыз.",
        "invalidMode" to "Көрсету режимін таңдаңыз.", "invalidLanguage" to "Тілді таңдаңыз.",
        "passwordShort" to "Кемінде 12 таңба пайдаланыңыз.", "passwordMismatch" to "Құпиясөздер сәйкес келмейді.",
        "incorrectPassword" to "Құпиясөз қате. Қайталап көріңіз.", "saveFailed" to "Өзгерістер сақталмады. Құрылғы жадын тексеріп, қайталаңыз.",
        "loadFailed" to "Сайт жүктелмеді. Байланысты тексеріңіз немесе баптаулардан мекенжайды өзгертіңіз.",
        "configFailed" to "Баптаулар оқылмады. Құрылғы жадын тексеріңіз немесе Sitewrap қолданбасын қайта орнатыңыз.",
        "extensionFailed" to "Бет масштабы қосылмады. Sitewrap қолданбасын қайта іске қосыңыз немесе орнатыңыз.",
        "unsupportedLink" to "Бұл сілтеме үшін басқа қолданба немесе веб-мекенжай қажет."
    )
    private val it = mapOf(
        "settings" to "Impostazioni", "unlock" to "Sblocca le impostazioni", "password" to "Password delle impostazioni", "open" to "Sblocca",
        "url" to "Indirizzo del sito", "zoom" to "Zoom della pagina (%)", "mode" to "Modalità di visualizzazione", "language" to "Lingua",
        "newPassword" to "Nuova password", "confirm" to "Conferma password", "save" to "Salva e apri", "cancel" to "Annulla",
        "window" to "Finestra", "kiosk" to "Chiosco", "version" to "Versione", "retry" to "Riprova",
        "invalidUrl" to "Inserisci un indirizzo HTTP(S) completo senza credenziali.", "invalidZoom" to "Scegli uno zoom tra il 30% e il 500%.",
        "invalidMode" to "Scegli una modalità di visualizzazione.", "invalidLanguage" to "Scegli una lingua.",
        "passwordShort" to "Usa almeno 12 caratteri.", "passwordMismatch" to "Le password non coincidono.",
        "incorrectPassword" to "Password errata. Riprova.", "saveFailed" to "Modifiche non salvate. Controlla la memoria del dispositivo e riprova.",
        "loadFailed" to "Impossibile caricare il sito. Controlla la connessione o modifica l'indirizzo nelle impostazioni.",
        "configFailed" to "Impossibile leggere le impostazioni. Controlla la memoria o reinstalla Sitewrap.",
        "extensionFailed" to "Impossibile attivare lo zoom. Riavvia o reinstalla Sitewrap.",
        "unsupportedLink" to "Questo link richiede un'altra app o un indirizzo web."
    )
    fun get(language: String, key: String) = mapOf("de" to de, "kk" to kk, "it" to it)[language]?.get(key) ?: en[key] ?: key
}
