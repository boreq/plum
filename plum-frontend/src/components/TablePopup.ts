import { defineComponent, type PropType } from 'vue';
import type { TableHeader, TableRow } from '@/dto/Table';
import Table from '@/components/Table.vue';

export default defineComponent({
    name: 'TablePopup',

    components: {
        Table,
    },

    props: {
        title: {
            type: String,
            default: '',
        },
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
    },

    emits: ['close'],

    methods: {
        close(): void {
            this.$emit('close');
        },
    },
});
