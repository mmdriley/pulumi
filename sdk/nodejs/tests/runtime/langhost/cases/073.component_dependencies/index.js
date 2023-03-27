let assert = require("assert");
let pulumi = require("../../../../../");

class MyCustomResource extends pulumi.CustomResource {
	constructor(name, opts) {
		super("test:index:MyCustomResource", name, {}, opts);
	}
}

class MyComponentResource extends pulumi.ComponentResource {
	constructor(name, opts) {
		super("test:index:MyComponentResource", name, {}, opts);
	}
}

let first = new MyComponentResource("first");
let firstChild = new MyCustomResource("firstChild", {
	parent: first,
});
let second = new MyComponentResource("second", {
	parent: first,
	dependsOn: first,
});
let myresource = new MyCustomResource("myresource", {
	parent: second,
});
