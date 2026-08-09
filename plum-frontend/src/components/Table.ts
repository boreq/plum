import { defineComponent, type PropType } from 'vue';
import { SortDirection, type TableHeader, type TableRow, type TableValue } from '@/dto/Table';

export interface IndexedRow {
    row: TableRow;
    index: number;
}

export default defineComponent({
    name: 'Table',

    props: {
        header: {
            type: Object as PropType<TableHeader>,
            default: null,
        },
        rows: {
            type: Array as PropType<TableRow[]>,
            default: () => [],
        },
        perPage: {
            type: Number,
            default: 10,
        },
        clickable: {
            type: Boolean,
            default: false,
        },
        expandable: {
            type: Boolean,
            default: true,
        },
    },

    emits: ['click-row'],

    data() {
        return {
            page: 0,
            sortColumn: null as number,
            sortDirection: null as SortDirection,
            expanded: false,
            query: '',
        };
    },

    computed: {
        dataPresent(): boolean {
            return this.rows && this.rows.length > 0;
        },

        indexedRows(): IndexedRow[] {
            return this.rows.map((row, index) => ({ row, index }));
        },

        sortedRows(): IndexedRow[] {
            if (this.sortColumn === null || this.sortDirection === null) {
                return this.indexedRows;
            }

            const direction = this.sortDirection === SortDirection.Ascending ? 1 : -1;

            return Array.from(this.indexedRows).sort((a, b) => {
                const result = this.compare(
                    a.row.data[this.sortColumn],
                    b.row.data[this.sortColumn],
                );
                return result === 0 ? a.index - b.index : result * direction;
            });
        },

        limitedRows(): IndexedRow[] {
            const start = this.page * this.perPage;
            return this.sortedRows.slice(start, start + this.perPage);
        },

        matchingRows(): IndexedRow[] {
            const query = this.query.trim().toLowerCase();
            if (!query) {
                return this.sortedRows;
            }
            return this.sortedRows.filter(entry => {
                return entry.row.data.some((value, columnIndex) => {
                    return this.formatValue(columnIndex, value).toLowerCase().includes(query);
                });
            });
        },

        hasNextPage(): boolean {
            return this.page < this.lastPage;
        },

        hasPrevPage(): boolean {
            return this.page > 0;
        },

        lastPage(): number {
            if (!this.rows || this.rows.length === 0) {
                return 0;
            }
            return Math.ceil(this.rows.length / this.perPage) - 1;
        },
    },

    watch: {
        rows(): void {
            if (this.page > this.lastPage) {
                this.page = this.lastPage;
            }
        },
    },

    mounted(): void {
        window.addEventListener('keydown', this.keydown);
    },

    beforeUnmount(): void {
        window.removeEventListener('keydown', this.keydown);
        this.unlockScroll();
    },

    methods: {
        compare(a: TableValue, b: TableValue): number {
            if (typeof a === 'number' && typeof b === 'number') {
                return a - b;
            }
            return String(a).localeCompare(String(b), undefined, { numeric: true });
        },

        formatValue(columnIndex: number, value: TableValue): string {
            const column = this.header.columns[columnIndex];
            return column.format ? column.format(value) : String(value);
        },

        isSortable(columnIndex: number): boolean {
            return this.header.columns[columnIndex].sortable !== false;
        },

        toggleSort(columnIndex: number): void {
            if (!this.isSortable(columnIndex)) {
                return;
            }

            if (this.sortColumn !== columnIndex) {
                this.sortColumn = columnIndex;
                this.sortDirection = SortDirection.Descending;
            } else if (this.sortDirection === SortDirection.Descending) {
                this.sortDirection = SortDirection.Ascending;
            } else {
                this.sortColumn = null;
                this.sortDirection = null;
            }

            this.page = 0;
        },

        sortIconClass(): string {
            return this.sortDirection === SortDirection.Ascending ? 'fas fa-sort-up' : 'fas fa-sort-down';
        },

        getColumnStyle(columnIndex: number): string {
            const column = this.header.columns[columnIndex];
            const styles: string[] = [];

            if (column.width) {
                styles.push(`width: ${column.width}`);
            } else {
                styles.push('flex: 1');
            }

            if (column.align) {
                styles.push(`text-align: ${column.align}`);
            }

            return styles.join(';');
        },

        getBackgroundStyle(row: TableRow): string {
            if (row.fraction) {
                const percentage = Math.round(row.fraction * 100);
                return `width: ${percentage}%;`;
            }
            return 'display: none;';
        },

        prevPage(): void {
            if (this.hasPrevPage) {
                this.page -= 1;
            }
        },

        nextPage(): void {
            if (this.hasNextPage) {
                this.page += 1;
            }
        },

        click(index: number): void {
            this.$emit('click-row', index);
        },

        clickExpanded(index: number): void {
            this.$emit('click-row', index);
            this.collapse();
        },

        expand(): void {
            this.expanded = true;
            this.query = '';
            this.lockScroll();
            this.$nextTick(() => {
                const search = this.$refs.search as HTMLInputElement;
                if (search) {
                    search.focus();
                }
            });
        },

        collapse(): void {
            this.expanded = false;
            this.unlockScroll();
        },

        keydown(event: KeyboardEvent): void {
            if (event.key === 'Escape' && this.expanded) {
                this.collapse();
            }
        },

        lockScroll(): void {
            document.body.style.overflow = 'hidden';
        },

        unlockScroll(): void {
            document.body.style.overflow = '';
        },
    },
});
