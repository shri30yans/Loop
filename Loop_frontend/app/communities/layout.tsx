import { getAllProjects } from '../projects/actions';

export default async function CommunitiesLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const projects = await getAllProjects();

  return (
    <div>
      {/* Use projects data here */}
      {children}
    </div>
  )
}
