export enum Browser {
    Chrome = 'chrome',
    Chromium = 'chromium',
    Firefox = 'firefox',
    Safari = 'safari',
}

const genericBrowserIcon = 'fas fa-window-maximize';

const browserIcons: { [browser in Browser]: string } = {
    [Browser.Chrome]: 'fab fa-chrome',
    [Browser.Chromium]: 'fab fa-chrome',
    [Browser.Firefox]: 'fab fa-firefox',
    [Browser.Safari]: 'fab fa-safari',
};

export function browserIcon(browser: string): string {
    const icon = browserIcons[browser as Browser];
    return icon ? icon : genericBrowserIcon;
}
