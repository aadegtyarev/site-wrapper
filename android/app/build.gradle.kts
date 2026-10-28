plugins { id("com.android.application") }

android {
    namespace = "io.github.aadegtyarev.sitewrap"
    compileSdk { version = release(37) { minorApiLevel = 2 } }
    defaultConfig {
        applicationId = "io.github.aadegtyarev.sitewrap"
        minSdk = 26
        targetSdk = 36
        versionCode = 1
        versionName = "0.1.0"
    }
    splits { abi {
        isEnable = true
        reset()
        include("arm64-v8a", "armeabi-v7a", "x86_64")
        isUniversalApk = false
    } }
    compileOptions { sourceCompatibility = JavaVersion.VERSION_17; targetCompatibility = JavaVersion.VERSION_17 }
}

dependencies {
    implementation("androidx.core:core-ktx:1.19.1")
    implementation("org.mozilla.geckoview:geckoview:159.0.20261022150525")
    testImplementation("junit:junit:4.13.2")
}
