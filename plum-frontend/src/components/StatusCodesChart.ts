import { defineComponent, markRaw, type PropType } from 'vue';
import type { Dictionary, SeriesPoint } from '@/dto/Data';
import { ChartColors, dimmed } from '@/dto/ChartColors';
import { TextService } from '@/services/TextService';
import { GroupingType } from '@/dto/GroupingType';
import { ChartAnimation } from '@/dto/ChartAnimation';
import { Chart } from '@/chartjs';

class ChartData {
    label: string;
    statuses: Dictionary<number>;
}

type BarChart = Chart<'bar', number[], string>;

const STATUS_COLORS = [
    ChartColors.Blue,
    ChartColors.Green,
    ChartColors.Violet,
    ChartColors.Orange,
    ChartColors.Red,
];

const textService = new TextService();

export default defineComponent({
    name: 'StatusCodesChart',

    props: {
        data: {
            type: Array as PropType<SeriesPoint[]>,
            default: () => [],
        },
        selectedIndex: {
            type: Number as PropType<number>,
            default: null,
        },
        groupingType: {
            type: Number as PropType<GroupingType>,
            default: GroupingType.Hourly,
        },
    },

    emits: ['select-data'],

    data() {
        return {
            chart: null as BarChart,
        };
    },

    watch: {
        data(): void {
            this.redraw();
        },
        selectedIndex(): void {
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

            const chartData: ChartData[] = this.data.map(point => {
                return {
                    label: textService.formatDate(point.time, this.groupingType),
                    statuses: point.statuses || {},
                };
            });
            this.drawChart(chartData);
        },

        toStatusString(status: string): string {
            return status[0] + 'xx';
        },

        drawChart(chartData: ChartData[]): void {
            const labels: string[] = chartData.map(v => v.label);

            if (!this.chart) {
                this.chart = markRaw(this.createChart());
            }

            this.chart.data.labels.length = 0;
            for (const label of labels) {
                this.chart.data.labels.push(label);
            }

            const statusDatas: number[][] = [];
            for (let datasetIndex = 0; datasetIndex < this.chart.data.datasets.length; datasetIndex++) {
                const statusType = this.toStatusString((datasetIndex + 1).toString());
                statusDatas.push(chartData.map(v => {
                    const total = Object.values(v.statuses).reduce((acc, hits) => acc + hits, 0);
                    return total ? (v.statuses[statusType] || 0) / total : 0;
                }));
            }

            for (let datasetIndex = 0; datasetIndex < this.chart.data.datasets.length; datasetIndex++) {
                this.chart.data.datasets[datasetIndex].data.length = statusDatas[datasetIndex].length;
                for (let dataIndex = 0; dataIndex < statusDatas[datasetIndex].length; dataIndex++) {
                    if (!this.chart.data.datasets[datasetIndex].data[dataIndex]) {
                        this.chart.data.datasets[datasetIndex].data[dataIndex] = 0;
                    }
                }
            }
            this.chart.update('none');

            for (let datasetIndex = 0; datasetIndex < this.chart.data.datasets.length; datasetIndex++) {
                statusDatas[datasetIndex].forEach((value, statusIndex) => {
                    this.chart.data.datasets[datasetIndex].data[statusIndex] = value;
                });
                const color = STATUS_COLORS[datasetIndex];
                this.chart.data.datasets[datasetIndex].backgroundColor = statusDatas[datasetIndex].map((_, index) => {
                    return this.selectedIndex === null || this.selectedIndex === index ? color : dimmed(color);
                });
            }
            this.chart.update();
        },

        createChart(): BarChart {
            return new Chart<'bar', number[], string>(this.$refs.canvas as HTMLCanvasElement, {
                type: 'bar',
                data: {
                    labels: [],
                    datasets: STATUS_COLORS.map((color, index) => {
                        return this.createDataset(this.toStatusString((index + 1).toString()), color);
                    }),
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
                            min: 0,
                            max: 1,
                            ticks: {
                                maxTicksLimit: 5,
                                callback: value => {
                                    return Math.round(Number(value) * 100) + '%';
                                },
                            },
                            stacked: true,
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
                                    const label = context.dataset.label;
                                    const percent = Math.round(context.parsed.y * 100);
                                    return `${label}: ${percent}%`;
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
                backgroundColor: color,
                borderWidth: 0,
                barPercentage: 1,
                categoryPercentage: 0.9,
            };
        },
    },
});
