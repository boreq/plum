export class UriService {

    private readonly staticResourceExtensions: string[] = [
        '.css',
        '.js',

        '.ico',
        '.png',
        '.jpg',
        '.jpeg',
        '.svg',

        '.webm',

        '.xml',
        '.rss',
        '.txt',
        '.abe',

        '.otf',
        '.ttf',
        '.woff',
        '.woff2',
        '.eot',
    ];

    isStaticResource(uri: string): boolean {
        for (const staticResourceExtension of this.staticResourceExtensions) {
            const cleanedUpUri = this.removeQueryFromUri(uri);
            if (cleanedUpUri.endsWith(staticResourceExtension)) {
                return true;
            }
        }
        return false;
    }

    private removeQueryFromUri(uri: string): string {
        return uri.split('?')[0];
    }

}
