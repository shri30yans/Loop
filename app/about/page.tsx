import { title,sectionheading} from "@/components/primitives";

export default function AboutPage() {
  return (
    <div className="space-y-5 *:space-y-2">

      <div>
          <div>
            <h1 className={title()}>About</h1>
          </div>
          <div>
              Welcome to DreamForge, the ultimate platform where creativity meets collaboration. At DreamForge, we believe that every groundbreaking innovation starts with a dream. 
              Our mission is to provide a vibrant community and powerful tools to help visionaries, creators, and learners turn their ideas into reality.
          </div>
      </div>

      <div>
        <h1 className={sectionheading({size : "lg"})}>Mission</h1>
        <p>
          Our mission is to empower individuals and teams to bring their creative visions to life by offering a collaborative, resource-rich environment. Whether you're an aspiring entrepreneur, a seasoned innovator, or someone with a passion for creating, DreamForge is designed to support you every step of the way.
        </p>
      </div>

      <div>
        <h1 className={sectionheading({size : "lg"})}>Features</h1>
        <ul className="space-y-4 ml-4">
          <li>
            <h1 className="text-2xl font-semibold">Community and Collaboration:</h1>
            <p>DreamForge is a place where you can connect with like-minded individuals from around the world. Engage in meaningful discussions, collaborate on projects, and build lasting relationships. Our platform fosters a supportive community that thrives on shared knowledge and collective growth.</p>
          </li>
          <li>
            <h1 className="text-2xl font-semibold">Project Management Tools:</h1>
            <p>Stay organized and on track with our suite of project management tools. From milestone tracking to visual progress reports, DreamForge provides everything you need to manage your projects efficiently. Our intuitive tools help you stay accountable and make continuous progress.</p>
          </li>
          <li>
            <h1 className="text-2xl font-semibold">Resource Repository:</h1>
            <p>Access a vast library of resources, including research papers, guides, templates, and tools tailored to your needs. Our curated repository ensures that you have the best materials at your fingertips to support your project development.</p>
          </li>
          <li>
            <h1 className="text-2xl font-semibold">Mentorship and Feedback:</h1>
            <p>Benefit from our peer review system and mentorship programs. Receive constructive feedback from experienced creators and industry experts, helping you refine your ideas and achieve your goals.</p>
          </li>
          <li>
            <h1 className="text-2xl font-semibold">AI-Powered Assistance:</h1>
            <p>Leverage the power of artificial intelligence with our intelligent project assistants. From idea generation and brainstorming to contextual recommendations and automated summaries, our AI tools enhance your creative process and provide valuable insights.</p>
          </li>
          <li>
            <h1 className="text-2xl font-semibold">Blogs:</h1>
            <p>Explore our collection of insightful blogs written by industry experts and thought leaders. Stay up-to-date with the latest trends, best practices, and success stories in your field. Our blogs provide valuable insights and inspiration to fuel your creative journey.</p>
          </li>
        </ul>
      </div>
      
      <div>
        <h1 className={sectionheading({size : "lg"})}>Vision</h1>
        <div>
          We envision DreamForge as a global hub for innovation and creativity. Our platform is designed to break down barriers, foster cross-cultural collaboration, and democratize access to resources and mentorship. By bringing together a diverse community of dreamers and doers, we aim to drive the next wave of innovation.
          Join Us
          Whether you're at the beginning of your creative journey or looking to take your projects to the next level, DreamForge is here to support you. Join our community and start forging your dreams into reality today.
        </div>
      </div>
    </div>
    
  );
}
