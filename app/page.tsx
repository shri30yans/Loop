import { Link } from "@nextui-org/link";
import { Snippet } from "@nextui-org/snippet";
import { Code } from "@nextui-org/code";
import { button as buttonStyles } from "@nextui-org/theme";

import { siteConfig } from "@/config/site";
import { title, subtitle } from "@/components/primitives";
import { GithubIcon } from "@/components/icons";

export default function Home() {
  return (
    <section className="flex flex-col items-center justify-center gap-4 py-8 md:py-10">
      <div className="inline-block max-w-3xl text-center justify-center">
        <h1 className={title({size: "lg"})}>Create with the </h1>
        <h1 className={title({size: "lg", color: "blue" })}>world.</h1>
        <br />
        <h2 className={subtitle({size: "md", class: "mt-4" })}>
          The best way to build your next big thing, find inspiration, and share your progress.
        </h2>
      </div>
    </section>
  );
}
