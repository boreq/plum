export enum Browser {
    Chrome = 'chrome',
    Chromium = 'chromium',
    Firefox = 'firefox',
    Safari = 'safari',
    Edge = 'edge',
    Opera = 'opera',
    Brave = 'brave',
    Vivaldi = 'vivaldi',
    DuckDuckGo = 'duckduckgo',
    SamsungInternet = 'samsung-internet',
    Yandex = 'yandex-browser',
    InternetExplorer = 'internet-explorer',
}

const genericBrowserIcon = 'fas fa-window-maximize';

const browserIcons: { [browser in Browser]: string } = {
    [Browser.Chrome]: 'fab fa-chrome',
    [Browser.Chromium]: 'fab fa-chrome',
    [Browser.Firefox]: 'fab fa-firefox',
    [Browser.Safari]: 'fab fa-safari',
    [Browser.Edge]: 'fab fa-edge',
    [Browser.Opera]: 'fab fa-opera',
    [Browser.Brave]: genericBrowserIcon,
    [Browser.Vivaldi]: genericBrowserIcon,
    [Browser.DuckDuckGo]: genericBrowserIcon,
    [Browser.SamsungInternet]: genericBrowserIcon,
    [Browser.Yandex]: 'fab fa-yandex',
    [Browser.InternetExplorer]: 'fab fa-internet-explorer',
};

export function browserIcon(browser: string): string {
    const icon = browserIcons[browser as Browser];
    return icon ? icon : genericBrowserIcon;
}
