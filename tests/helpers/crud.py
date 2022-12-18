from observer.context import Context
from observer.entities.projects import NewProject, Project


async def create_project(ctx: Context, name: str, description: str) -> Project:
    project = await ctx.projects_repo.create_project(NewProject(name=name, description=description))
    return project
