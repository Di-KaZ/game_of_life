export interface Grid {
	width: number;
	height: number;
	paused: boolean;
	cells: Cell[]
}

export interface Cell { alive: boolean, turns: number }
