export enum FilterDimension {
    Category = 'category',
    Uri = 'uri',
    Status = 'status',
    Referer = 'referer',
}

export type Filter = {
    [dimension in FilterDimension]?: string;
};

export const FilterLabels: { [dimension in FilterDimension]: string } = {
    [FilterDimension.Category]: 'Category',
    [FilterDimension.Uri]: 'Page',
    [FilterDimension.Status]: 'Status',
    [FilterDimension.Referer]: 'Referer',
};

export function filtersEqual(a: Filter, b: Filter): boolean {
    return Object.values(FilterDimension).every(dimension => a[dimension] === b[dimension]);
}
