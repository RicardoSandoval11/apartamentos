using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using scheduler_service.DataAccess.Entities;

namespace scheduler_service.DataAccess.Config;

public class ReservationConfiguration : IEntityTypeConfiguration<Reservation>
{
    public void Configure(EntityTypeBuilder<Reservation> builder)
    {
        builder.ToTable("reservation");

        builder.Property(e => e.Id)
        .HasColumnName("id")
        .IsRequired();

        builder.HasKey(e => e.Id);

        builder.Property(e => e.TimeslotId)
        .HasColumnName("timeslot_id")
        .IsRequired();

        builder.Property(e => e.ReservationStatus)
        .HasColumnName("reservation_status")
        .IsRequired();

        builder.HasOne(e => e.Timeslot)
        .WithMany(e => e.Reservations)
        .HasForeignKey(e => e.TimeslotId);
    }
}