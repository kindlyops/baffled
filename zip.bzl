def _pkg_zip_impl(ctx):
    files = " ".join(["%s" % dep.path for dep in ctx.files.deps])
    cmd = """
rm -f {out}
{zipper} cf {out} {path}
"""
    cmd = cmd.format(
        out = ctx.outputs.zip.path,
        path = files,
        zipper = ctx.executable._zipper.path,
    )

    ctx.actions.run_shell(
        inputs = ctx.files.deps + [ctx.executable._zipper],
        outputs = [ctx.outputs.zip],
        command = cmd,
        progress_message = "making pkg_zip archive %s (%s files)" % (ctx.label, len(files)),
    )

pkg_zip = rule(
    implementation = _pkg_zip_impl,
    attrs = {
        "deps": attr.label_list(allow_files = True),
        "_zipper": attr.label(
            default = Label("@bazel_tools//tools/zip:zipper"),
            allow_files = True,
            executable = True,
            cfg = "host",
        ),
    },
    outputs = {
        "zip": "pkg_%{name}",
    },
)
