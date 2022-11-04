// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import assert from "assert";

import {
    fullyQualifiedStackName,
    isFullyQualifiedStackName,
    LocalWorkspace,
    RemoteGitAuthArgs,
    RemoteGitProgramArgs,
    RemoteStack,
    RemoteWorkspace,
    RemoteWorkspaceOptions,
} from "../../automation";
import { asyncTest } from "../util";
import { getTestOrg, getTestSuffix } from "./util";

describe("isFullyQualifiedStackName", () => {
    const tests = [
        {
            name: "fully qualified",
            input: "owner/project/stack",
            expected: true,
        },
        {
            name: "undefined",
            input: undefined,
            expected: false,
        },
        {
            name: "null",
            input: null,
            expected: false,
        },
        {
            name: "empty",
            input: "",
            expected: false,
        },
        {
            name: "name",
            input: "name",
            expected: false,
        },
        {
            name: "name & owner",
            input: "owner/name",
            expected: false,
        },
        {
            name: "sep",
            input: "/",
            expected: false,
        },
        {
            name: "two seps",
            input: "//",
            expected: false,
        },
        {
            name: "three seps",
            input: "///",
            expected: false,
        },
        {
            name: "invalid",
            input: "owner/project/stack/wat",
            expected: false,
        },
    ];

    tests.forEach(test => {
        it(`${test.name}`, () => {
            const actual = isFullyQualifiedStackName(test.input!);
            assert.strictEqual(actual, test.expected);
        });
    });
});

describe("RemoteWorkspace", () => {
    describe("selectStack", () => {
        describe("throws appropriate errors", () => testErrors(RemoteWorkspace.selectStack));
    });
    describe("createStack", () => {
        describe("throws appropriate errors", () => testErrors(RemoteWorkspace.createStack));
        it(`runs through the stack lifecycle`, asyncTest(testLifecycle(RemoteWorkspace.createStack)));
    });
    describe("createOrSelectStack", () => {
        describe("throws appropriate errors", () => testErrors(RemoteWorkspace.createOrSelectStack));
        xit(`runs through the stack lifecycle`, asyncTest(testLifecycle(RemoteWorkspace.createOrSelectStack)));
    });
});

const testRepo = "https://github.com/pulumi/test-repo.git";

function testErrors(fn: (args: RemoteGitProgramArgs, opts?: RemoteWorkspaceOptions) => Promise<RemoteStack>) {
    const stack = "owner/project/stack";
    const tests: {
        name: string;
        stackName: string;
        url: string;
        branch?: string;
        commitHash?: string;
        auth?: RemoteGitAuthArgs;
        error: string;
    }[] = [
        {
            name: "stack empty",
            stackName: "",
            url: "",
            error: `"" stack name must be fully qualified`,
        },
        {
            name: "stack just name",
            stackName: "name",
            url: "",
            error: `"name" stack name must be fully qualified`,
        },
        {
            name: "stack just name & owner",
            stackName: "owner/name",
            url: "",
            error: `"owner/name" stack name must be fully qualified`,
        },
        {
            name: "stack just sep",
            stackName: "/",
            url: "",
            error: `"/" stack name must be fully qualified`,
        },
        {
            name: "stack just two seps",
            stackName: "//",
            url: "",
            error: `"//" stack name must be fully qualified`,
        },
        {
            name: "stack just three seps",
            stackName: "///",
            url: "",
            error: `"///" stack name must be fully qualified`,
        },
        {
            name: "stack invalid",
            stackName: "owner/project/stack/wat",
            url: "",
            error: `"owner/project/stack/wat" stack name must be fully qualified`,
        },
        {
            name: "no url",
            stackName: stack,
            url: "",
            error: `url is required.`,
        },
        {
            name: "no branch or commit",
            stackName: stack,
            url: testRepo,
            error: `at least commitHash or branch are required.`,
        },
        {
            name: "both branch and commit",
            stackName: stack,
            url: testRepo,
            branch: "branch",
            commitHash: "commit",
            error: `commitHash and branch cannot both be specified.`,
        },
        {
            name: "both ssh private key and path",
            stackName: stack,
            url: testRepo,
            branch: "branch",
            auth: {
                sshPrivateKey: "key",
                sshPrivateKeyPath: "path",
            },
            error: `sshPrivateKey and sshPrivateKeyPath cannot both be specified.`,
        },
    ];

    tests.forEach(test => {
        it(`${test.name}`, asyncTest(async () => {
            const { stackName, url, branch, commitHash, auth } = test;

            let succeeded = false;
            let message: string | undefined = undefined;
            try{
                await fn({ stackName, url, branch, commitHash, auth });
                succeeded = true;
            } catch (err) {
                message = (<any>err).message;
            }
            assert.strictEqual(succeeded, false);
            assert.strictEqual(message, test.error);
        }));
    });
}

function testLifecycle(fn: (args: RemoteGitProgramArgs, opts?: RemoteWorkspaceOptions) => Promise<RemoteStack>) {
    return async () => {
        const stackName = fullyQualifiedStackName(getTestOrg(), "go_remote_proj", `int_test${getTestSuffix()}`);
        const stack = await fn({
            stackName,
            url: testRepo,
            branch: "refs/heads/master",
            projectPath: "goproj",
        },{
            preRunCommands: [
                `pulumi config set bar abc --stack ${stackName}`,
                `pulumi config set --secret buzz secret --stack ${stackName}`,
            ],
        });

        // pulumi up
        const upRes = await stack.up();
        assert.strictEqual(Object.keys(upRes.outputs).length, 3);
        assert.strictEqual(upRes.outputs["exp_static"].value, "foo");
        assert.strictEqual(upRes.outputs["exp_static"].secret, false);
        assert.strictEqual(upRes.outputs["exp_cfg"].value, "abc");
        assert.strictEqual(upRes.outputs["exp_cfg"].secret, false);
        assert.strictEqual(upRes.outputs["exp_secret"].value, "secret");
        assert.strictEqual(upRes.outputs["exp_secret"].secret, true);
        assert.strictEqual(upRes.summary.kind, "update");
        assert.strictEqual(upRes.summary.result, "succeeded");

        // pulumi preview
        const preRes = await stack.preview();
        assert.strictEqual(preRes.changeSummary.same, 1);

        // pulumi refresh
        const refRes = await stack.refresh();
        assert.strictEqual(refRes.summary.kind, "refresh");
        assert.strictEqual(refRes.summary.result, "succeeded");

        // pulumi destroy
        const destroyRes = await stack.destroy();
        assert.strictEqual(destroyRes.summary.kind, "destroy");
        assert.strictEqual(destroyRes.summary.result, "succeeded");

        await (await LocalWorkspace.create({})).removeStack(stackName);
    };
}
