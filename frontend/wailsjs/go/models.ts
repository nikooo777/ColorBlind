export namespace main {
	
	export class Defaults {
	    confuserColor: string;
	    shapeMode: string;
	    density: number;
	    minDensity: number;
	    maxDensity: number;
	    textDensity: number;
	    shade: number;
	    minShade: number;
	    maxShade: number;
	
	    static createFrom(source: any = {}) {
	        return new Defaults(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.confuserColor = source["confuserColor"];
	        this.shapeMode = source["shapeMode"];
	        this.density = source["density"];
	        this.minDensity = source["minDensity"];
	        this.maxDensity = source["maxDensity"];
	        this.textDensity = source["textDensity"];
	        this.shade = source["shade"];
	        this.minShade = source["minShade"];
	        this.maxShade = source["maxShade"];
	    }
	}
	export class PlateRequest {
	    confuserColor: string;
	    shapeMode: string;
	    shapeText: string;
	    density: number;
	    shade: number;
	    showOutline: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PlateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.confuserColor = source["confuserColor"];
	        this.shapeMode = source["shapeMode"];
	        this.shapeText = source["shapeText"];
	        this.density = source["density"];
	        this.shade = source["shade"];
	        this.showOutline = source["showOutline"];
	    }
	}
	export class PlateResponse {
	    image: string;
	    message: string;
	    placed: number;
	    requested: number;
	
	    static createFrom(source: any = {}) {
	        return new PlateResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.image = source["image"];
	        this.message = source["message"];
	        this.placed = source["placed"];
	        this.requested = source["requested"];
	    }
	}

}

