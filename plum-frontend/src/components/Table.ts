import { defineComponent, type PropType } from 'vue';
import type { TableHeader, TableRow } from '@/dto/Table';

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
    },

    emits: ['click-row'],

    data() {
        return {
            page: 0,
        };
    },

    computed: {
        dataPresent(): boolean {
            return this.rows && this.rows.length > 0;
        },

        limitedRows(): TableRow[] {
            const start = this.page * this.perPage;
            return this.rows.slice(start, start + this.perPage);
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

    methods: {
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

        getBackgroundStyle(rowIndex: number): string {
            const row = this.limitedRows[rowIndex];
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

        click(rowIndex: number): void {
            const index = this.perPage * this.page + rowIndex;
            this.$emit('click-row', index);
        },
    },
});
