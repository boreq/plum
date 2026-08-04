import { defineComponent, markRaw, type PropType } from 'vue';
import type { RangeData } from '@/dto/Data';
import { ChartColors } from '@/dto/ChartColors';
import { TextService } from '@/services/TextService';
import { GroupingType } from '@/dto/GroupingType';
import { ChartAnimation } from '@/dto/ChartAnimation';
import { Chart } from '@/chartjs';

class ChartData {
    label: string;
    bytes: number;
}

type LineChart = Chart<'line', number[], string>;

const textService = new TextService();

export default defineComponent({
    name: 'BytesSentChart',

    props: {
        data: {
            type: Array as PropType<RangeData[]>,
            default: () => [],
        },
        groupingType: {
            type: Number as PropType<GroupingType>,
            default: GroupingType.Hourly,
        },
    },

    emits: ['select-data'],

    data() {
        return {
            chart: null as LineChart,
        };
    },

    watch: {
        data(): void {
            this.redraw();
        },
    },

    mounted(): void {
        this.redraw();
    },

    beforeUnmount(): void {
        if (this.chart) {
            this.chart.destroy();
            this.chart = null;
        }
    },

    methods: {
        redraw(): void {
            if (!this.data) {
                return;
            }

            const chartData: ChartData[] = this.data
                .map(rangeData => {
                    return {
                        label: textService.formatDate(rangeData.time, this.groupingType),
                        bytes: rangeData.data.bytes,
                    };
                });
            this.drawChart(chartData);
        },

        drawChart(chartData: ChartData[]): void {
            const bytes: number[] = chartData.map(v => v.bytes);
            const labels: string[] = chartData.map(v => v.label);

            if (!this.chart) {
                this.chart = markRaw(this.createChart());
            }

            this.chart.data.labels.length = 0;
            for (const label of labels) {
                this.chart.data.labels.push(label);
            }

            this.chart.data.datasets[0].data.length = bytes.length;
            bytes.forEach((value, i) => {
                if (!this.chart.data.datasets[0].data[i]) {
                    this.chart.data.datasets[0].data[i] = 0;
                }
            });
            this.chart.update('none');

            bytes.forEach((value, i) => {
                this.chart.data.datasets[0].data[i] = value;
            });
            this.chart.update();
        },

        createChart(): LineChart {
            return new Chart<'line', number[], string>(this.$refs.canvas as HTMLCanvasElement, {
                type: 'line',
                data: {
                    labels: [],
                    datasets: [
                        this.createDataset('Bytes sent', ChartColors.Primary),
                    ],
                },
                options: {
                    maintainAspectRatio: false,
                    animation: {
                        duration: ChartAnimation.Duration,
                    },
                    scales: {
                        x: {
                            ticks: {
                                display: false,
                            },
                            grid: {
                                display: false,
                            },
                            stacked: true,
                        },
                        y: {
                            beginAtZero: true,
                            ticks: {
                                maxTicksLimit: 5,
                                callback: value => {
                                    return textService.humanizeBytes(Number(value));
                                },
                            },
                            stacked: false,
                        },
                    },
                    plugins: {
                        legend: {
                            display: false,
                        },
                        tooltip: {
                            mode: 'index',
                            callbacks: {
                                label: context => {
                                    return textService.humanizeBytes(context.parsed.y);
                                },
                            },
                        },
                    },
                    onClick: (_event, elements) => {
                        const index = elements.length > 0 ? elements[0].index : null;
                        this.$emit('select-data', index);
                    },
                },
            });
        },

        createDataset(label: string, color: string) {
            return {
                label: label,
                data: [] as number[],
                backgroundColor: 'rgba(0, 0, 0, 0)',
                pointBackgroundColor: color,
                borderColor: color,
                borderWidth: 2,
            };
        },
    },
});
