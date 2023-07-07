export namespace controller {
	
	export class ConfigItem {
	    id: number;
	    name: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new ConfigItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.value = source["value"];
	    }
	}

}

