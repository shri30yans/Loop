import connection from '@/utils/db'; 

// Function to fetch projects from the database
export async function fetchProjects(type: string, sortBy: string, timeRange: string) {
  try {
    let query = `
      SELECT *
      FROM PROJECT
      ORDER BY ${sortBy} -- Sort by the given sortBy field
    `;
    const [rows] = await connection.execute(query, [type]);
    
    return rows; 

  } catch (error) {
    console.error('Error fetching projects: ', error);
    return null; 
  }
}
