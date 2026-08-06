export enum Category {
    Malicious = 'malicious',
    Automated = 'automated',
    Unclassified = 'unclassified',
}

export const CategoryLabels: { [category in Category]: string } = {
    [Category.Malicious]: 'Malicious traffic',
    [Category.Automated]: 'Automated traffic',
    [Category.Unclassified]: 'Unclassified traffic',
};
