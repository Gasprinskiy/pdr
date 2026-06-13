export namespace dialogs {
	
	export class OpenFileDialogParam {
	    FileExtensions: string[];
	    LangCode: number;
	
	    static createFrom(source: any = {}) {
	        return new OpenFileDialogParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FileExtensions = source["FileExtensions"];
	        this.LangCode = source["LangCode"];
	    }
	}

}

