import { getAllProjects } from '../projects/actions';

export default async function CommunitiesLayout({
  children,
}: {
  children: React.ReactNode
}) {

  return (
    <div>
      {/* Use projects data here */}
      {children}
    </div>
  )
}
