# Site Wrapper

Небольшое приложение для запуска заданного сайта на движке Firefox/Gecko.

> Статус: прототип. Desktop-код и Android-приложение находятся в разработке. Артефакты CI нельзя считать релизными до подписи расширения и проверки на Linux, Windows и Android.

## Требования

- Linux и Windows: запуск как отдельного приложения; URL из настроек или первый URL-аргумент командной строки.
- Android: отдельный APK на GeckoView; URL из настроек, а внешний URL — через Android Intent.
- Масштаб страницы и пароль на изменение настроек.
- Язык интерфейса настроек и предпочтительный язык сайта: английский (по умолчанию), немецкий, казахский, итальянский. Версия приложения показана в настройках.
- Два режима отображения: отдельное окно приложения и полноэкранный киоск.
- Локальная разработка и сборки в Docker; затем сборки для трёх платформ в CI репозитория.

## Архитектура

### Desktop

Тонкий нативный запускатор + установленный Firefox с отдельным профилем + собственное WebExtension. В оконном режиме расширение открывает сайт в окне типа `popup` без адресной строки и вкладок и управляет масштабом. В режиме киоска запускатор открывает Firefox с `--kiosk`, расширение применяет тот же масштаб. Запускатор хранит настройки, проверяет пароль и передаёт конфигурацию расширению через Native Messaging. Исходники Firefox собирать не нужно. Если Firefox не найден, приложение предлагает установить его и показывает ссылку на официальный сайт Mozilla.

Сочетание `--kiosk` с расширением и поведение новых вкладок нужно проверить в прототипе. Режим киоска скрывает интерфейс браузера, но сам по себе не блокирует доступ к ОС.

### Android

Небольшое Kotlin-приложение с GeckoView. Экран настроек защищён паролем; экран сайта использует те же URL и масштаб. Оконному режиму соответствует обычный экран приложения, киоску — полноэкранный immersive-режим. Запуск по URL реализуется через Intent, поскольку аргументов командной строки у обычного Android-приложения нет. Если понадобится запрет выхода из приложения, это отдельная функция Android Lock Task с настройкой устройства.

### Общие правила

- Приоритет URL: явный URL запуска, затем сохранённый URL.
- Принимаем только `http://` и `https://`; ошибки показываем до запуска браузера.
- Пароль храним в виде стойкого хеша с солью. Он защищает интерфейс настроек, но не является защитой от пользователя с доступом к файлам приложения.
- Отдельный профиль Firefox не смешивает cookies, историю и расширения с личным браузером пользователя.

## Сборка и проверка

Первый релиз описан в [контракте](docs/contracts/features/initial-release.md). Отдельные контракты задают [UI](docs/contracts/surfaces/ui.md), [UX](docs/contracts/surfaces/ux.md), [экран настроек](docs/contracts/surfaces/settings-screen.md) и [качество кода](docs/contracts/surfaces/code-quality.md).

Локальная проверка без установки Go в систему:

```sh
make check
make desktop-package
make android-check
```

`make desktop-package` складывает исполняемый файл Linux и неподписанное расширение в `dist/`. Их нужно держать рядом. XPI отправляется в Mozilla Add-ons на [подпись в режиме `unlisted`](https://extensionworkshop.com/documentation/publish/signing-and-distribution-overview/), без публикации в каталоге; полученный подписанный XPI распространяется вместе с приложением. Неподписанный XPI предназначен для разработки и не запустится в стабильной версии Firefox. Для Windows CI собирает отдельный `.exe`. Артефакты CI содержат неподписанное расширение.

Для Android CI собирает отладочные APK отдельно для `arm64-v8a`, `armeabi-v7a` и `x86_64`. Выбирайте APK под архитектуру устройства; общий APK с тремя архитектурами был около 509 МБ. `make android-check` загружает SDK и Gradle в Docker-тома, запускает тесты, Lint и сборку, не устанавливая инструменты в систему. Первый запуск принимает лицензии Android SDK. GeckoView берётся как готовая зависимость Mozilla, браузер из исходников не собирается.

Desktop: `site-wrapper [URL]`, `site-wrapper run [URL]`, `site-wrapper version`. Первый запуск открывает настройку, затем сайт. Настройки открываются `Ctrl+Alt+S` или кнопкой в оконном режиме. Android принимает ссылку через системное действие «Открыть с помощью». В киоске Android настройки открываются удержанием двух пальцев 1,8 секунды.

Перед публикацией нужно подписать desktop XPI, проверить запуск и переходы на Linux и Windows, а Android — на реальном устройстве. Масштаб Android применяется CSS `zoom` через встроенное расширение GeckoView; он требует проверки на целевых сайтах, поскольку ведёт себя иначе, чем масштаб вкладки Firefox на desktop.

## Первоисточники

- [Firefox WebExtensions: окна без панелей](https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/windows/create)
- [Firefox WebExtensions: масштаб](https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/tabs/setZoom)
- [Firefox WebExtensions: Native Messaging](https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Native_messaging)
- [Firefox: параметры командной строки](https://firefox-source-docs.mozilla.org/browser/CommandLineParameters.html)
- [GeckoView](https://firefox-source-docs.mozilla.org/mobile/android/geckoview/index.html)

## Товарный знак

Firefox is a trademark of the Mozilla Foundation in the U.S. and other countries. This project is independent and is not affiliated with Mozilla. Название Site Wrapper не использует товарный знак Firefox; зависимость от браузера указана отдельно.
