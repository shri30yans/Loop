import { Link } from "@nextui-org/link";
import { Snippet } from "@nextui-org/snippet";
import { Code } from "@nextui-org/code";
import { button as buttonStyles } from "@nextui-org/theme";
import { AuthProvider } from "@/components/auth/AuthProvider";

import { siteConfig } from "@/config/site";
import { landingpagetext, subtitle, subheading } from "@/components/ui/primitives";
import { GithubIcon } from "@/components/ui/icons";

export default function Home() {
  return (
    <section className="flex flex-col items-center justify-center">
    <div className="flex h-screen pb-11" >
      <div className="inline-block max-w-3xl text-center justify-center m-auto">
        <h1 className={landingpagetext({size: "lg"})}>Create with the </h1>
        <h1 className={landingpagetext({size: "lg", color: "blue" })}>world.</h1>
        <h2 className={subtitle({size: "lg", class: "mt-4" })}>
          Build the next big thing.
        </h2>
      </div>
    </div>
        
        <div className="space-y- *:space-y-12 text-3xl max-w-5xl items-center text-center justify-center gap-4 py-8 md:py-10">
          {/* <div>
              <div >
                  The ultimate platform where creativity meets collaboration. At Loop, we believe that every groundbreaking innovation starts with a dream. 
                  Our mission is to provide a vibrant community and powerful tools to help visionaries, creators, and learners turn their ideas into reality.
              </div>
          </div> */}

          <div>
            <ul className="space-y-12 *:space-y-2">
              <li>
                <h1 className="text-5xl font-semibold">Community and Collaboration:</h1>
                <p>Loop is a place where you can connect with like-minded individuals from around the world. Engage in meaningful discussions, collaborate on projects, and build lasting relationships.</p>
              </li>
              <li>
                <h1 className="text-5xl font-semibold">Project Management</h1>
                <p>Stay organized and on track with progress tracking and recieving feedback on every step of the way.</p>
              </li>
              <li>
                <h1 className="text-5xl font-semibold">Build in public</h1>
                <p>Garner an audience and find people interested in your project before it's launch.</p>
              </li>
              <li>
                <h1 className="text-5xl font-semibold">Mentorship and Feedback</h1>
                <p>Benefit from our peer review system and mentorship programs. Receive constructive feedback from experienced creators and industry experts, helping you refine your ideas and achieve your goals.</p>
              </li>
              <li>
                <h1 className="text-5xl font-semibold">AI-Powered Assistance</h1>
                <p>Leverage the power of artificial intelligence from brainstorming to contextual recommendations, ai assisted editing and automated summaries, our AI tools enhance your creative process and provide valuable insights.</p>
              </li>
            </ul>
          </div>
        
    
      </div>
    </section>
  );
}
