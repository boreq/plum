import { defineComponent, markRaw, type PropType } from 'vue';
import type { SeriesPoint } from '@/dto/Data';
import { ChartColors, dimmed } from '@/dto/ChartColors';
import { TextService } from '@/services/TextService';
import { GroupingType } from '@/dto/GroupingType';
import { ChartAnimation } from '@/dto/ChartAnimation';
import { Chart } from '@/chartjs';

class ChartData {
    label: string;
    hits: number;
    visits: number;
}

type BarChart = Chart<'bar', number[], string>;

const textService = new TextService();

export default defineComponent({
    name: 'HitsAndVisits',

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

            const chartData: ChartData[] = this.data
                .map(point => {
                    return {
                        label: textService.formatDate(point.time, this.groupingType),
                        hits: point.hits,
                        visits: point.visits,
                    };
                });
            this.drawChart(chartData);
        },

        drawChart(chartData: ChartData[]): void {
            const visits: number[] = chartData.map(v => v.visits);
            const hits: number[] = chartData.map(v => v.hits);
            const labels: string[] = chartData.map(v => v.label);

            if (!this.chart) {
                this.chart = markRaw(this.createChart());
            }

            this.chart.data.labels.length = 0;
            for (const label of labels) {
                this.chart.data.labels.push(label);
            }

            this.chart.data.datasets[0].data.length = visits.length;
            visits.forEach((value, index) => {
                if (!this.chart.data.datasets[0].data[index]) {
                    this.chart.data.datasets[0].data[index] = 0;
                }
            });

            this.chart.data.datasets[1].data.length = hits.length;
            hits.forEach((value, index) => {
                if (!this.chart.data.datasets[1].data[index]) {
                    this.chart.data.datasets[1].data[index] = 0;
                }
            });

            this.chart.update('none');

            visits.forEach((value, index) => {
                this.chart.data.datasets[0].data[index] = value;
            });

            hits.forEach((value, index) => {
                this.chart.data.datasets[1].data[index] = value;
            });

            this.chart.data.datasets[0].backgroundColor = this.barColors(ChartColors.Primary, visits.length);
            this.chart.data.datasets[1].backgroundColor = this.barColors(ChartColors.Secondary, hits.length);

            this.chart.update();
        },

        barColors(color: string, count: number): string[] {
            return Array.from({length: count}, (_, index) => {
                return this.selectedIndex === null || this.selectedIndex === index ? color : dimmed(color);
            });
        },

        createChart(): BarChart {
            return new Chart<'bar', number[], string>(this.$refs.canvas as HTMLCanvasElement, {
                type: 'bar',
                data: {
                    labels: [],
                    datasets: [
                        this.createDataset('Visits', ChartColors.Primary),
                        this.createDataset('Hits', ChartColors.Secondary),
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
                                    const n = Number(value);
                                    if (n >= 10000) {
                                        return n / 1000 + 'k';
                                    }
                                    return value;
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
