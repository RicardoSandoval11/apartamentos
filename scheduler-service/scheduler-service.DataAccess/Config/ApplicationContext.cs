using Microsoft.EntityFrameworkCore;

namespace scheduler_service.DataAccess.Config;

public class ApplicationContext(DbContextOptions<ApplicationContext> options) : DbContext(options)
{
    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.ApplyConfiguration(new ReservationConfiguration());
        modelBuilder.ApplyConfiguration(new TimeslotConfiguration());
        modelBuilder.ApplyConfiguration(new ReservingConfiguration());
        base.OnModelCreating(modelBuilder);
    }
}