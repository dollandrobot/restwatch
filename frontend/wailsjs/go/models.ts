export namespace main {
	
	export class Message {
	    id: string;
	    // Go type: time
	    receivedAt: any;
	    method: string;
	    body: string;
	    bodyMarkdown: string;
	    contentLength: number;
	    remoteAddr: string;
	    header: Record<string, string[]>;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.receivedAt = this.convertValues(source["receivedAt"], null);
	        this.method = source["method"];
	        this.body = source["body"];
	        this.bodyMarkdown = source["bodyMarkdown"];
	        this.contentLength = source["contentLength"];
	        this.remoteAddr = source["remoteAddr"];
	        this.header = source["header"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UserOptions {
	    port: number;
	    maxMessagesToKeep: number;
	    defaultEndpoint: string;
	    jumpToLatest: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UserOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.port = source["port"];
	        this.maxMessagesToKeep = source["maxMessagesToKeep"];
	        this.defaultEndpoint = source["defaultEndpoint"];
	        this.jumpToLatest = source["jumpToLatest"];
	    }
	}

}

