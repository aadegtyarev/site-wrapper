package io.github.aadegtyarev.sitewrap

import android.app.Activity
import android.content.Intent
import android.graphics.Color
import android.os.Bundle
import android.os.Handler
import android.os.Looper
import android.text.InputType
import android.util.Log
import android.view.Gravity
import android.view.MotionEvent
import android.view.View
import android.view.WindowManager
import android.widget.Button
import android.widget.EditText
import android.widget.FrameLayout
import android.widget.LinearLayout
import android.widget.ScrollView
import android.widget.Spinner
import android.widget.ArrayAdapter
import android.widget.AdapterView
import android.widget.TextView
import org.json.JSONObject
import org.mozilla.geckoview.GeckoResult
import org.mozilla.geckoview.GeckoRuntime
import org.mozilla.geckoview.GeckoRuntimeSettings
import org.mozilla.geckoview.GeckoSession
import org.mozilla.geckoview.GeckoView
import org.mozilla.geckoview.WebExtension
import org.mozilla.geckoview.WebRequestError

class MainActivity : Activity() {
    companion object { private var sharedRuntime: GeckoRuntime? = null }
    private lateinit var store: ConfigStore
    private lateinit var root: FrameLayout
    private var config: Config? = null
    private var launchUrl: String? = null
    private var session: GeckoSession? = null
    private var currentUrl: String? = null
    private val handler = Handler(Looper.getMainLooper())
    private var gesture: Runnable? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        store = ConfigStore(this)
        config = store.load()
        root = FrameLayout(this)
        setContentView(root)
        launchUrl = intent.dataString?.takeIf { candidate -> Config(candidate, 100, "window", "en").error() == null }
        if (intent.dataString != null && launchUrl == null) {
            Log.w("Sitewrap", "Invalid external URL")
            showError("invalidUrl") { if (config == null) showSettings(true) else showSite() }
            return
        }
        when {
            store.hasConfig() && config == null -> showError("configFailed")
            config == null -> showSettings(true)
            else -> showSite()
        }
    }

    override fun onNewIntent(intent: Intent) {
        super.onNewIntent(intent)
        val candidate = intent.dataString ?: return
        if (Config(candidate, 100, "window", "en").error() != null) {
            showError("invalidUrl") { if (config == null) showSettings(true) else showSite() }
            return
        }
        launchUrl = candidate
        currentUrl = candidate
        if (config != null) showSite()
    }

    override fun dispatchTouchEvent(event: MotionEvent): Boolean {
        if (config?.mode == "kiosk") {
            if (event.pointerCount == 2 && gesture == null) {
                gesture = Runnable { gesture = null; showUnlock() }.also { handler.postDelayed(it, 1800) }
            } else if (event.pointerCount != 2 || event.actionMasked == MotionEvent.ACTION_UP || event.actionMasked == MotionEvent.ACTION_CANCEL) {
                gesture?.let(handler::removeCallbacks)
                gesture = null
            }
        }
        return super.dispatchTouchEvent(event)
    }

    private fun word(key: String) = Words.get(config?.language ?: "en", key)
    private fun dp(value: Int) = (value * resources.displayMetrics.density).toInt()
    private fun field(label: String, password: Boolean = false): EditText {
        val input = EditText(this)
        input.hint = label
        input.contentDescription = label
        input.textSize = 16f
        input.inputType = if (password) InputType.TYPE_CLASS_TEXT or InputType.TYPE_TEXT_VARIATION_PASSWORD else InputType.TYPE_CLASS_TEXT
        input.minHeight = dp(52)
        return input
    }
    private fun layout(title: String): LinearLayout {
        window.decorView.systemUiVisibility = 0
        root.removeAllViews()
        val content = LinearLayout(this).apply {
            orientation = LinearLayout.VERTICAL
            setPadding(dp(24), dp(24), dp(24), dp(24))
            setBackgroundColor(Color.rgb(248, 250, 252))
        }
        val heading = TextView(this).apply { text = title; textSize = 26f; setTextColor(Color.rgb(22, 36, 56)); setPadding(0, 0, 0, dp(20)) }
        content.addView(heading)
        val scroll = ScrollView(this).apply { isFillViewport = true; addView(content) }
        root.addView(scroll)
        return content
    }
    private fun label(content: LinearLayout, value: String): TextView = TextView(this).also {
        it.text = value; it.textSize = 16f; it.setTextColor(Color.rgb(49, 65, 83)); it.setPadding(0, dp(12), 0, dp(4)); content.addView(it)
    }
    private fun button(content: LinearLayout, title: String, action: () -> Unit) = Button(this).also {
        it.text = title; it.minHeight = dp(52); it.setOnClickListener { action() }; content.addView(it)
    }
    private fun showError(key: String, retry: (() -> Unit)? = null) {
        val content = layout(word("settings"))
        label(content, word(key))
        retry?.let { button(content, word("retry"), it) }
        if (config != null) button(content, word("settings")) { showUnlock() }
    }
    private fun showUnlock() {
        val content = layout(word("unlock"))
        val password = field(word("password"), true).also(content::addView)
        val error = label(content, "")
        button(content, word("open")) {
            try {
                if (store.verify(password.text.toString())) showSettings(false)
                else { error.text = word("incorrectPassword"); password.text.clear() }
            } catch (failure: Exception) {
                Log.e("Sitewrap", "Password verifier failed", failure)
                error.text = word("configFailed")
            }
        }
        button(content, word("cancel")) { showSite() }
    }
    private fun showSettings(first: Boolean) {
        val content = layout(word("settings"))
        val saved = config ?: Config("", 100, "window", "en")
        val url = field(word("url")).also { it.setText(saved.url); content.addView(it) }
        val zoom = field(word("zoom")).also { it.inputType = InputType.TYPE_CLASS_NUMBER; it.setText(saved.zoom.toString()); content.addView(it) }
        val modeLabel = label(content, word("mode"))
        val mode = Spinner(this).also {
            it.adapter = ArrayAdapter(this, android.R.layout.simple_spinner_dropdown_item, listOf(word("window"), word("kiosk")))
            it.setSelection(if (saved.mode == "kiosk") 1 else 0); content.addView(it)
        }
        val languageLabel = label(content, word("language"))
        val language = Spinner(this).also {
            it.adapter = ArrayAdapter(this, android.R.layout.simple_spinner_dropdown_item, listOf("English", "Deutsch", "Қазақша", "Italiano"))
            it.setSelection(listOf("en", "de", "kk", "it").indexOf(saved.language).coerceAtLeast(0)); content.addView(it)
        }
        val password = field(word("newPassword"), true).also(content::addView)
        val confirmation = field(word("confirm"), true).also(content::addView)
        val version = label(content, "${word("version")}: ${packageManager.getPackageInfo(packageName, 0).versionName}")
        val error = label(content, "")
        val save = button(content, word("save")) {
            val draft = Config(url.text.toString().trim(), zoom.text.toString().toIntOrNull() ?: -1,
                if (mode.selectedItemPosition == 1) "kiosk" else "window", listOf("en", "de", "kk", "it")[language.selectedItemPosition])
            val failure = draft.error() ?: when {
                first && password.text.length < 12 -> "passwordShort"
                password.text.isNotEmpty() && password.text.length < 12 -> "passwordShort"
                password.text.toString() != confirmation.text.toString() -> "passwordMismatch"
                else -> null
            }
            if (failure != null) {
                val message = Words.get(draft.language, failure)
                when (failure) {
                    "invalidUrl" -> url.error = message
                    "invalidZoom" -> zoom.error = message
                    "passwordShort" -> password.error = message
                    "passwordMismatch" -> confirmation.error = message
                    else -> error.text = message
                }
                return@button
            }
            try {
                if (!store.save(draft, password.text.toString().takeIf { it.isNotEmpty() })) {
                    error.text = word("saveFailed"); return@button
                }
                config = draft
                if (launchUrl == null) currentUrl = draft.url
                showSite()
            } catch (failure: Exception) {
                Log.e("Sitewrap", "Settings save failed", failure)
                error.text = word("saveFailed")
            }
        }
        val cancel = if (!first) button(content, word("cancel")) { showSite() } else null
        language.onItemSelectedListener = object : AdapterView.OnItemSelectedListener {
            override fun onNothingSelected(parent: AdapterView<*>?) = Unit
            override fun onItemSelected(parent: AdapterView<*>?, view: View?, position: Int, id: Long) {
                val selected = listOf("en", "de", "kk", "it")[position]
                fun t(key: String) = Words.get(selected, key)
                (content.getChildAt(0) as TextView).text = t("settings")
                url.hint = t("url"); url.contentDescription = t("url")
                zoom.hint = t("zoom"); zoom.contentDescription = t("zoom")
                modeLabel.text = t("mode"); languageLabel.text = t("language")
                password.hint = t("newPassword"); password.contentDescription = t("newPassword")
                confirmation.hint = t("confirm"); confirmation.contentDescription = t("confirm")
                version.text = "${t("version")}: ${packageManager.getPackageInfo(packageName, 0).versionName}"
                save.text = t("save"); cancel?.text = t("cancel")
                val selectedMode = mode.selectedItemPosition
                mode.adapter = ArrayAdapter(this@MainActivity, android.R.layout.simple_spinner_dropdown_item, listOf(t("window"), t("kiosk")))
                mode.setSelection(selectedMode)
            }
        }
    }
    private fun showSite() {
        val saved = config ?: return
        session?.close()
        session = null
        root.removeAllViews()
        window.decorView.systemUiVisibility = if (saved.mode == "kiosk") 5894 or View.SYSTEM_UI_FLAG_IMMERSIVE_STICKY else 0
        if (saved.mode == "kiosk") window.addFlags(WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON) else window.clearFlags(WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON)
        val view = GeckoView(this)
        root.addView(view, FrameLayout.LayoutParams(-1, -1))
        if (saved.mode == "window") {
            val settings = Button(this).apply { text = "⚙"; contentDescription = word("settings"); setOnClickListener { showUnlock() } }
            root.addView(settings, FrameLayout.LayoutParams(dp(52), dp(52), Gravity.TOP or Gravity.END))
        }
        val runtime = sharedRuntime ?: GeckoRuntime.create(applicationContext, GeckoRuntimeSettings.Builder().locales(arrayOf(saved.language)).build()).also { sharedRuntime = it }
        runtime.settings.setLocales(arrayOf(saved.language))
        val opened = GeckoSession()
        session = opened
        opened.setNavigationDelegate(object : GeckoSession.NavigationDelegate {
            override fun onLoadError(session: GeckoSession, uri: String?, error: WebRequestError): GeckoResult<String>? {
                Log.e("Sitewrap", "Site load failed: ${error.code}")
                runOnUiThread { showError("loadFailed") { showSite() } }
                return null
            }
        })
        opened.setProgressDelegate(object : GeckoSession.ProgressDelegate {
            override fun onPageStart(session: GeckoSession, url: String) {
                if (url.startsWith("http://") || url.startsWith("https://")) currentUrl = url
            }
        })
        opened.open(runtime)
        view.setSession(opened)
        runtime.webExtensionController.ensureBuiltIn("resource://android/assets/zoom/", "android-zoom@sitewrap").accept({ extension: WebExtension? ->
            runOnUiThread {
                if (this.session !== opened) return@runOnUiThread
                if (extension == null) showError("extensionFailed") { showSite() }
                else {
                    opened.webExtensionController.setMessageDelegate(extension, object : WebExtension.MessageDelegate {
                        override fun onMessage(nativeApp: String, message: Any, sender: WebExtension.MessageSender): GeckoResult<Any> {
                            val percent = if ((message as? JSONObject)?.optString("action") == "zoom") saved.zoom else 100
                            return GeckoResult.fromValue(JSONObject().put("percent", percent) as Any)
                        }
                    }, "sitewrap")
                    opened.loadUri(currentUrl ?: launchUrl ?: saved.url)
                }
            }
        }, { failure: Throwable? ->
            Log.e("Sitewrap", "Zoom extension failed", failure)
            runOnUiThread { showError("extensionFailed") { showSite() } }
        })
    }

    override fun onDestroy() {
        session?.close()
        super.onDestroy()
    }
}
