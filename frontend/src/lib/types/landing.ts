export interface ERPModule {
	id: string;
	code: string;
	name: string;
	category: 'Keuangan' | 'Rantai Pasok' | 'Penjualan & CRM' | 'Operasional' | 'SDM' | 'Sistem';
	description: string;
	icon: string;
	color: string;
	badge?: string;
	features: string[];
	inputsFrom?: string[];
	outputsTo?: string[];
}

export interface MetricItem {
	value: string;
	label: string;
	sublabel: string;
	icon: string;
}

export interface SolutionItem {
	title: string;
	tagline: string;
	description: string;
	modules: string[];
	icon: string;
	highlights: string[];
}
