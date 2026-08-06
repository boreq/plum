export enum Category {
    Automated = 'automated',
    Unclassified = 'unclassified',
}

export const CategoryLabels: { [category in Category]: string } = {
    [Category.Automated]: 'Automated traffic',
    [Category.Unclassified]: 'Unclassified traffic',
};
