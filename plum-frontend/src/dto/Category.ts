export enum Category {
    Malicious = 'malicious',
    Automated = 'automated',
    PossiblyAutomated = 'possibly-automated',
    Unclassified = 'unclassified',
}

export const CategoryLabels: { [category in Category]: string } = {
    [Category.Malicious]: 'Malicious',
    [Category.Automated]: 'Automated',
    [Category.PossiblyAutomated]: 'Possibly automated',
    [Category.Unclassified]: 'Unclassified',
};
