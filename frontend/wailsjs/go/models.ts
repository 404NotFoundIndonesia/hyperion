export namespace main {
	
	export class ObfuscationConfig {
	    RemoveWhiteSpace: boolean;
	    RemoveComments: boolean;
	    ChangeVariable: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ObfuscationConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RemoveWhiteSpace = source["RemoveWhiteSpace"];
	        this.RemoveComments = source["RemoveComments"];
	        this.ChangeVariable = source["ChangeVariable"];
	    }
	}

}

