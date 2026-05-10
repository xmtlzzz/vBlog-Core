export namespace proto {
	
	export class TrendPoint {
	    label?: string;
	    pv?: number;
	    uv?: number;
	    view_total?: number;
	    comment_count?: number;
	    post_count?: number;
	
	    static createFrom(source: any = {}) {
	        return new TrendPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.label = source["label"];
	        this.pv = source["pv"];
	        this.uv = source["uv"];
	        this.view_total = source["view_total"];
	        this.comment_count = source["comment_count"];
	        this.post_count = source["post_count"];
	    }
	}
	export class GetTrendsResponse {
	    points?: TrendPoint[];
	
	    static createFrom(source: any = {}) {
	        return new GetTrendsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.points = this.convertValues(source["points"], TrendPoint);
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
	export class LatestStats {
	    pv_today?: number;
	    uv_today?: number;
	    total_posts?: number;
	    total_views?: number;
	    total_comments?: number;
	    total_tags?: number;
	    pv_yesterday?: number;
	    uv_yesterday?: number;
	    views_today_delta?: number;
	    comments_today_delta?: number;
	
	    static createFrom(source: any = {}) {
	        return new LatestStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pv_today = source["pv_today"];
	        this.uv_today = source["uv_today"];
	        this.total_posts = source["total_posts"];
	        this.total_views = source["total_views"];
	        this.total_comments = source["total_comments"];
	        this.total_tags = source["total_tags"];
	        this.pv_yesterday = source["pv_yesterday"];
	        this.uv_yesterday = source["uv_yesterday"];
	        this.views_today_delta = source["views_today_delta"];
	        this.comments_today_delta = source["comments_today_delta"];
	    }
	}

}

