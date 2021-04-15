import {
  PieChartType,
  SingleStatBaseOptions,
  PieChartLabels,
  PieChartLegendOptions,
  GraphTooltipOptions,
} from '@grafana/ui';
export interface PieChartOptions extends SingleStatBaseOptions {
  pieType: PieChartType;
  displayLabels: PieChartLabels[];
  legend: PieChartLegendOptions;
  tooltip: GraphTooltipOptions;
}
