import { defineComponent, markRaw, type PropType } from 'vue';
import type { Data, Dictionary, RangeData } from '@/dto/Data';
import { ChartColors } from '@/dto/ChartColors';
import { TextService } from '@/services/TextService';
import { GroupingType } from '@/dto/GroupingType';
import { ChartAnimation } from '@/dto/ChartAnimation';
import { Chart } from '@/chartjs';

class ChartData {
    label: string;
    statuses: Dictionary<number>;
}

type BarChart = Chart<'bar', number[], string>;

const textService = new TextService();

export default defineComponent({
    name: 'StatusCodesChart',

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
            chart: null as BarChart,
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

            const chartData: ChartData[] = this.data.map(rangeData => {
                return {
                    label: textService.formatDate(rangeData.time, this.groupingType),
                    statuses: this.groupByStatusType(rangeData.data),
                };
            });
            this.drawChart(chartData);
        },

        groupByStatusType(data: Data): Dictionary<number> {
            const rv: Dictionary<number> = {};
            if (!data || !data.statuses) {
                return rv;
            }
            Object.entries(data.statuses).forEach(([status, metrics]) => {
                const statusType = this.toStatusString(status);
                rv[statusType] = (rv[statusType] || 0) + metrics.hits;
            });
            return rv;
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
            }
            this.chart.update();
        },

        createChart(): BarChart {
            return new Chart<'bar', number[], string>(this.$refs.canvas as HTMLCanvasElement, {
                type: 'bar',
                data: {
                    labels: [],
                    datasets: [
                        this.createDataset('1xx', ChartColors.Blue),
                        this.createDataset('2xx', ChartColors.Green),
                        this.createDataset('3xx', ChartColors.Violet),
                        this.createDataset('4xx', ChartColors.Orange),
                        this.createDataset('5xx', ChartColors.Red),
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
                pointBackgroundColor: color,
                borderColor: color,
                borderWidth: 2,
                barPercentage: 1,
                categoryPercentage: 0.9,
            };
        },
    },
});
